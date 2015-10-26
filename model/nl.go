package model

/*
#cgo CFLAGS: -DNO_REUSE
#include "asl.h"
*/
import "C"

import (
	"math"
	"unsafe"
)

type Problem struct {
	Name	string
	asl	*C.struct_ASL
}

func (p *Problem) NumVariables() int {
	return int(p.asl.i.n_var_)
}

func (p *Problem) Constraints() []Constraint {
	numConstraints := int(p.asl.i.n_con_)
	numNonLinear := int(p.asl.i.nlc_) 
	numNonLinearNetwork := int(p.asl.i.nlnc_) 
	numLinearNetwork := int(p.asl.i.lnc_)
	constraints := make([]Constraint, numConstraints)
	bounds := (*[1 << 30]C.real)(unsafe.Pointer(p.asl.i.LUrhs_))[:numConstraints*2:numConstraints*2]
	for i := 0; i < numConstraints; i++ {
		name := C.GoString(C.con_name_ASL(p.asl, C.int(i)))

		/* Get the constraint shape from position within the list of constraints */ 
		var conShape ConstraintShape
		if i < numNonLinear - numNonLinearNetwork {
			conShape = NonLinearGeneralConstraint
		} else if i < numNonLinear { 
			conShape = NonLinearNetworkConstraint
		} else if i < numConstraints - numLinearNetwork {
			conShape = LinearNetworkConstraint
		} else {
			conShape = LinearGeneralConstraint
		}

		/* Get the constraint type */
		upperIsInf := math.IsInf(float64(bounds[i*2]), 1)
		lowerIsInf := math.IsInf(float64(bounds[i*2+1]), 0)
		var conType ConstraintType
		if upperIsInf && !lowerIsInf {
			conType = ConstraintGreaterThan	
		} else if !upperIsInf && lowerIsInf {
			conType = ConstraintLessThan
		} else if bounds[i*2] == bounds[i*2+1] {
			conType = ConstraintEqualTo	
		} else if !upperIsInf && !lowerIsInf {
			conType = ConstraintRange
		} else {
			conType = ConstraintNonBinding		}
		constraints[i].Name = name
		constraints[i].Shape = conShape
		constraints[i].Type = conType
		constraints[i].Min = float64(bounds[i*2])
		constraints[i].Max = float64(bounds[i*2+1])
	}
	return constraints
}

func (p *Problem) Objectives() []Objective {
	numObjectives := int(p.asl.i.n_obj_)
	
	objectives := make([]Objective, numObjectives)
	objectiveSenses := (*[1<<30]byte)(unsafe.Pointer(p.asl.i.objtype_))[:numObjectives:numObjectives]
	for i := 0; i < numObjectives; i++ {
		name := C.GoString(C.obj_name_ASL(p.asl, C.int(i)))
		objectives[i].Name = name
		objectives[i].Sense = ObjectiveSense(objectiveSenses[i])
		// TODO: Objective shape
	}	
	
	return objectives
}

func (p *Problem) Variables() []Variable {
	numVariables := int(p.asl.i.n_var_)
	numLinearBinary := int(p.asl.i.nbv_)	
	numLinearNonBinaryInt := int(p.asl.i.niv_)
	numNonLinear := int(p.asl.i.nlvb_)
	numNetwork := int(p.asl.i.nwv_)	
	variables := make([]Variable, numVariables)
	bounds := (*[1 << 30]C.real)(unsafe.Pointer(p.asl.i.LUv_))[:numVariables*2:numVariables*2]
	for i := 0; i < numVariables; i++ {
		name := C.GoString(C.var_name_ASL(p.asl, C.int(i)))
		var varType VariableType
		variables[i].Name = name
		if i < numNonLinear {
			varType = VariableNonLinear
		} else if i < numNetwork + numNonLinear {
			varType = VariableNetworkLinear
		} else if i < numVariables - (numLinearBinary + numLinearNonBinaryInt) {
			varType = VariableOtherLinear

		} else if i < numVariables - numLinearBinary {
			varType = VariableLinearBinary

		} else {
			varType = VariableInteger
		}
		variables[i].Type = varType
		variables[i].LowerBound = float64(bounds[i*2])
		variables[i].UpperBound = float64(bounds[i*2+1])
	}
	return variables
}

/* Load a problem from a `.nl` file */
func ProblemFromFile(path string) (*Problem) {
	pathC := C.CString(path)
	asl := C.ASL_alloc(C.ASL_read_fg)
	nl := C.jac0dim_ASL(asl, pathC, C.ftnlen(len(path)))
	C.fg_read_ASL(asl, nl, 0)
	return &Problem {path, asl}
}
