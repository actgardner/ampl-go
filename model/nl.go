package model

/* A package for interacting with AMPL models. */

/*
#cgo CFLAGS: -DNO_REUSE
#include "asl.h"

// Helper functions to call gradient and value function pointers
// Go doesn't support C function pointers
typedef real (valFunc) (ASL* asl, int n, real *X, fint *nerror);
typedef void (grdFunc) (ASL* asl, int n, real *X, real *G, fint *nerror);
typedef void (jacFunc) (ASL* asl, real *X, real *J, fint *nerror);

real callValFunc(void *f, ASL *asl, int n, real *X, fint *nerror) {
	return ((valFunc*)f)(asl, n, X, nerror);
}

void callGrdFunc(void *f, ASL *asl, int n, real *X, real *G, fint *nerror) {
	return ((grdFunc*)f)(asl, n, X, G, nerror);
}

void callJacFunc(void *f, ASL *asl, real *X, real *J, fint *nerror) {
	return ((jacFunc*)f)(asl, X, J, nerror);
}

*/
import "C"

import (
	"fmt"
	"math"
	"unsafe"
)

type Problem struct {
	Name	string
	asl	*C.struct_ASL
}

/* Get the list of Constraints in this problem */
func (p *Problem) Constraints() []Constraint {
	numConstraints := int(p.asl.i.n_con_)
	numNonLinear := int(p.asl.i.nlc_) 
	numNonLinearNetwork := int(p.asl.i.nlnc_) 
	numLinearNetwork := int(p.asl.i.lnc_)
	constraints := make([]Constraint, numConstraints)
	bounds := (*[1 << 30]C.real)(unsafe.Pointer(p.asl.i.LUrhs_))[:numConstraints*2:numConstraints*2]
	cgradList := (*[1 << 30]*C.struct_cgrad)(unsafe.Pointer(p.asl.i.Cgrad_))[:numConstraints:numConstraints]
	vars := p.Variables()
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
		upperIsInf := math.IsInf(float64(bounds[i*2+1]), 1)
		lowerIsInf := math.IsInf(float64(bounds[i*2]), 0)

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
			conType = ConstraintNonBinding
		}

		constraints[i].Name = name
		constraints[i].Shape = conShape
		constraints[i].Type = conType
		constraints[i].Min = float64(bounds[i*2])
		constraints[i].Max = float64(bounds[i*2+1])
		constraints[i].Variables = make([]Gradient, 0)
		constraints[i].Index = i
		constraints[i].p = p
		gradPtr := cgradList[i]
		for gradPtr != nil {
			constraints[i].Variables = append(constraints[i].Variables, Gradient{vars[gradPtr.varno], float64(gradPtr.coef)})
			gradPtr = (*gradPtr).next
		}
	}
	return constraints
}

/* Get the list of Objectives in this problem */
func (p *Problem) Objectives() []Objective {
	numObjectives := int(p.asl.i.n_obj_)
	
	objectives := make([]Objective, numObjectives)
	objectiveSenses := (*[1<<30]byte)(unsafe.Pointer(p.asl.i.objtype_))[:numObjectives:numObjectives]
	ogradList := (*[1 << 30]*C.struct_ograd)(unsafe.Pointer(p.asl.i.Ograd_))[:numObjectives:numObjectives]
	vars := p.Variables()	
	for i := 0; i < numObjectives; i++ {
		name := C.GoString(C.obj_name_ASL(p.asl, C.int(i)))
		objectives[i].Name = name
		objectives[i].Sense = ObjectiveSense(objectiveSenses[i])
		objectives[i].Index = i
		objectives[i].p = p
	
		gradPtr := ogradList[i]
		for gradPtr != nil {
			objectives[i].Variables = append(objectives[i].Variables, Gradient{vars[gradPtr.varno], float64(gradPtr.coef)})
			gradPtr = (*gradPtr).next
		}
	}	
	
	return objectives
}

func (p *Problem) objValue(index int, x []float64) (float64, error) {
	numVariables := int(p.asl.i.n_var_)
	if len(x) != numVariables {
		return 0, fmt.Errorf("Error: Incorrect number of variables in input: expected %i, got %i", numVariables, len(x)) 
	}
	var err C.fint
	val := C.callValFunc(unsafe.Pointer(p.asl.p.Objval), p.asl, C.int(index), (*C.real)(unsafe.Pointer(&x[0])), &err)
	if err != 0 {
		return 0, fmt.Errorf("Error: %i when evaluating objective value", err)
	}
	return float64(val), nil
}

func (p *Problem) objGrad(index int, x []float64) ([]float64, error) {
	numVariables := int(p.asl.i.n_var_)
	if len(x) != numVariables {
		return nil, fmt.Errorf("Error: Incorrect number of variables in input: expected %i, got %i", numVariables, len(x)) 
	}
	var err C.fint
	grad := make([]float64, numVariables)
	C.callGrdFunc(unsafe.Pointer(p.asl.p.Objgrd), p.asl, C.int(index), (*C.real)(unsafe.Pointer(&x[0])), (*C.real)(unsafe.Pointer(&grad[0])), &err)
	if err != 0 {
		return nil, fmt.Errorf("Error: %i when evaluating objective value", err)
	}
	return grad, nil
}

func (p *Problem) conValue(index int, x []float64) (float64, error) {
	numVariables := int(p.asl.i.n_var_)
	if len(x) != numVariables {
		return 0, fmt.Errorf("Error: Incorrect number of variables in input: expected %i, got %i", numVariables, len(x)) 
	}
	var err C.fint
	val := C.callValFunc(unsafe.Pointer(p.asl.p.Conival), p.asl, C.int(index), (*C.real)(unsafe.Pointer(&x[0])), &err)
	if err != 0 {
		return 0, fmt.Errorf("Error: %i when evaluating constraint value", err)
	}
	return float64(val), nil
}

func (p *Problem) conGrad(index int, x []float64) ([]float64, error) {
	numVariables := int(p.asl.i.n_var_)
	if len(x) != numVariables {
		return nil, fmt.Errorf("Error: Incorrect number of variables in input: expected %i, got %i", numVariables, len(x)) 
	}
	var err C.fint
	grad := make([]float64, numVariables)
	C.callGrdFunc(unsafe.Pointer(p.asl.p.Congrd), p.asl, C.int(index), (*C.real)(unsafe.Pointer(&x[0])), (*C.real)(unsafe.Pointer(&grad[0])), &err)
	if err != 0 {
		return nil, fmt.Errorf("Error: %i when evaluating constraint value", err)
	}
	return grad, nil
}

/* Evaluate the value of all constraints at point x */
func (p *Problem) ConstraintValues(x []float64) ([]float64, error) {
	numVariables := int(p.asl.i.n_var_)
	numConstraints := int(p.asl.i.n_con_)
	if len(x) != numVariables {
		return nil, fmt.Errorf("Error: Incorrect number of variables in input: expected %i, got %i", numVariables, len(x)) 
	}
	var err C.fint
	vals := make([]float64, numConstraints)
	C.callJacFunc(unsafe.Pointer(p.asl.p.Conval), p.asl, (*C.real)(unsafe.Pointer(&x[0])), (*C.real)(unsafe.Pointer(&vals[0])), &err)
	if err != 0 {
		return nil, fmt.Errorf("Error: %i when evaluating constraint values", err)
	}
	return vals, nil
}

/* Evaluate the Jacobian of the constraints */
func (p *Problem) ConstraintJacobian(x []float64) ([]float64, error) {
	numVariables := int(p.asl.i.n_var_)
	numConstraints := int(p.asl.i.n_con_)
	if len(x) != numVariables {
		return nil, fmt.Errorf("Error: Incorrect number of variables in input: expected %i, got %i", numVariables, len(x)) 
	}
	var err C.fint
	vals := make([]float64, numConstraints * numVariables)
	C.callJacFunc(unsafe.Pointer(p.asl.p.Jacval), p.asl, (*C.real)(unsafe.Pointer(&x[0])), (*C.real)(unsafe.Pointer(&vals[0])), &err)
	if err != 0 {
		return nil, fmt.Errorf("Error: %i when evaluating constraint Jacobian", err)
	}
	return vals, nil
} 

/* Get the list of Variables in this problem */
func (p *Problem) Variables() []Variable {
	numVariables := int(p.asl.i.n_var_)
	numNonLinear := intMax(int(p.asl.i.nlvc_), int(p.asl.i.nlvo_))
	numBoth := int(p.asl.i.nlvb_)
	numBothInt := int(p.asl.i.nlvbi_)
	numConst := int(p.asl.i.nlvc_)
	numConstInt := int(p.asl.i.nlvci_)
	numObj := int(p.asl.i.nlvo_)
	numObjInt := int(p.asl.i.nlvoi_)
	numLinearArcs := int(p.asl.i.nwv_)	
	numBinary := int(p.asl.i.nbv_)	
	numNonBinaryInt := int(p.asl.i.niv_)

	variables := make([]Variable, numVariables)
	bounds := (*[1 << 30]C.real)(unsafe.Pointer(p.asl.i.LUv_))[:numVariables*2:numVariables*2]

	for i := 0; i < numVariables; i++ {
		name := C.GoString(C.var_name_ASL(p.asl, C.int(i)))
		var varType VariableType
		variables[i].Name = name
		if i < (numBoth - numBothInt) {	
			varType = VariableContinuousNonLinear
		} else if i < (numBoth) {
			varType = VariableIntegerNonLinear
		} else if i < (numConst - numConstInt) {
			varType = VariableContinuousNonLinear
		} else if i < numConst {
			varType = VariableIntegerNonLinear
		} else if i < (numConst + (numObj - (numBoth + numObjInt))) {
			varType = VariableContinuousNonLinear
		} else if i < numNonLinear {
			varType = VariableIntegerNonLinear
		} else if i < numLinearArcs + numNonLinear {
			varType = VariableLinearArc
		} else if i < numVariables - (numBinary + numNonBinaryInt) {
			varType = VariableOtherLinear
		} else if i < numVariables - numNonBinaryInt {
			varType = VariableBinary
		} else {
			varType = VariableOtherInteger
		}
		variables[i].Type = varType
		variables[i].LowerBound = float64(bounds[i*2])
		variables[i].UpperBound = float64(bounds[i*2+1])
		variables[i].Index = i
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

/* Return the larger integer */
func intMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}