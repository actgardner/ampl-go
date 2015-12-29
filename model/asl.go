package model

/* A package for interacting with AMPL models. */

/*
#cgo CFLAGS: -DNO_REUSE
#define PSHVREAD
#include "asl.h"
#include "psinfo.h"
#include "nlp2.h"

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
	aslPfgh	*C.struct_ASL_pfgh
}

/* Get the list of Constraints in this problem */
func (p *Problem) Constraints() []Constraint {
	numConstraints := int(p.asl.i.n_con_)
	constraints := make([]Constraint, numConstraints)
	bounds := (*[1 << 30]C.real)(unsafe.Pointer(p.asl.i.LUrhs_))[:numConstraints*2:numConstraints*2]
	cgradList := (*[1 << 30]*C.struct_cgrad)(unsafe.Pointer(p.asl.i.Cgrad_))[:numConstraints:numConstraints]
	cClassList := (*[1 << 30]C.char)(unsafe.Pointer(p.aslPfgh.I.c_class))[:numConstraints:numConstraints]
	vars := p.Variables()
	for i := 0; i < numConstraints; i++ {
		name := C.GoString(C.con_name_ASL(p.asl, C.int(i)))

		/* Get the constraint sense */
		upperIsInf := math.IsInf(float64(bounds[i*2+1]), 1)
		lowerIsInf := math.IsInf(float64(bounds[i*2]), 0)

		var conType ConstraintSense
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
		constraints[i].Sense = conType
		constraints[i].Min = float64(bounds[i*2])
		constraints[i].Max = float64(bounds[i*2+1])
		constraints[i].Variables = make([]Variable, 0)
		constraints[i].Index = i
		constraints[i].Shape = Shape(cClassList[i])
		constraints[i].p = p
		gradPtr := cgradList[i]
		for gradPtr != nil {
			constraints[i].Variables = append(constraints[i].Variables, vars[gradPtr.varno])
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
	oClassList := (*[1 << 30]C.char)(unsafe.Pointer(p.aslPfgh.I.o_class))[:numObjectives:numObjectives]
	vars := p.Variables()	
	for i := 0; i < numObjectives; i++ {
		name := C.GoString(C.obj_name_ASL(p.asl, C.int(i)))
		objectives[i].Name = name
		objectives[i].Sense = ObjectiveSense(objectiveSenses[i])
		objectives[i].Index = i
		objectives[i].Shape = Shape(oClassList[i])
		objectives[i].p = p
	
		gradPtr := ogradList[i]
		for gradPtr != nil {
			objectives[i].Variables = append(objectives[i].Variables, vars[gradPtr.varno])
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

	j := 0
	for i := 0; i < numBoth - numBothInt; i++ {
		name := C.GoString(C.var_name_ASL(p.asl, C.int(j)))
		variables[j].Name = name
		variables[j].Type = VariableReal
		variables[j].LowerBound = float64(bounds[j*2])
		variables[j].UpperBound = float64(bounds[j*2+1])
		variables[j].Index = j 
		j++
	}

	for i := 0; i < numBothInt; i++ {
		name := C.GoString(C.var_name_ASL(p.asl, C.int(j)))
		variables[j].Name = name
		variables[j].LowerBound = float64(bounds[j*2])
		variables[j].UpperBound = float64(bounds[j*2+1])
		// Integer variables may in fact be binary - check if the bounds are 0 and 1
		if variables[j].LowerBound == 0 && variables[j].UpperBound == 1 {
			variables[j].Type = VariableBinary
		} else {
			variables[j].Type = VariableInteger
		}
		variables[j].Index = j
		j++
	}

	for i := 0; i < numConst - (numBoth + numConstInt); i++ {
		name := C.GoString(C.var_name_ASL(p.asl, C.int(j)))
		variables[j].Name = name
		variables[j].Type = VariableReal
		variables[j].LowerBound = float64(bounds[j*2])
		variables[j].UpperBound = float64(bounds[j*2+1])
		variables[j].Index = j
		j++
	}

	for i := 0; i < numConstInt; i++ {
		name := C.GoString(C.var_name_ASL(p.asl, C.int(j)))
		variables[j].Name = name
		variables[j].LowerBound = float64(bounds[j*2])
		variables[j].UpperBound = float64(bounds[j*2+1])
		// Integer variables may in fact be binary - check if the bounds are 0 and 1
		if variables[j].LowerBound == 0 && variables[j].UpperBound == 1 {
			variables[j].Type = VariableBinary
		} else {
			variables[j].Type = VariableInteger
		}
		variables[j].Index = j
		j++
	}

	for i := 0; i < numObj - (numConst + numObjInt); i++ {
		name := C.GoString(C.var_name_ASL(p.asl, C.int(j)))
		variables[j].Name = name
		variables[j].Type = VariableReal
		variables[j].LowerBound = float64(bounds[j*2])
		variables[j].UpperBound = float64(bounds[j*2+1])
		variables[j].Index = j
		j++
	}

	for i := 0; i < numObjInt; i++ {
		name := C.GoString(C.var_name_ASL(p.asl, C.int(j)))
		variables[j].Name = name
		variables[j].LowerBound = float64(bounds[j*2])
		variables[j].UpperBound = float64(bounds[j*2+1])
		// Integer variables may in fact be binary - check if the bounds are 0 and 1
		if variables[j].LowerBound == 0 && variables[j].UpperBound == 1 {
			variables[j].Type = VariableBinary
		} else {
			variables[j].Type = VariableInteger
		}
		variables[j].Index = j
		j++
	}

	for i := 0; i < numLinearArcs; i++ {
		name := C.GoString(C.var_name_ASL(p.asl, C.int(j)))
		variables[j].Name = name
		variables[j].Type = VariableArc
		variables[j].LowerBound = float64(bounds[j*2])
		variables[j].UpperBound = float64(bounds[j*2+1])
		variables[j].Index = j
		j++
	}

	for i := 0; i < numVariables - (numNonLinear + numBinary + numNonBinaryInt + numLinearArcs); i++ {
		name := C.GoString(C.var_name_ASL(p.asl, C.int(j)))
		variables[j].Name = name
		variables[j].Type = VariableOther
		variables[j].LowerBound = float64(bounds[j*2])
		variables[j].UpperBound = float64(bounds[j*2+1])
		variables[j].Index = j
		j++
	}

	for i := 0; i < numBinary; i++ {
		name := C.GoString(C.var_name_ASL(p.asl, C.int(j)))
		variables[j].Name = name
		variables[j].Type = VariableBinary
		variables[j].LowerBound = float64(bounds[j*2])
		variables[j].UpperBound = float64(bounds[j*2+1])
		variables[j].Index = j
		j++
	}

	for i := 0; i < numNonBinaryInt; i++ {
		name := C.GoString(C.var_name_ASL(p.asl, C.int(j)))
		variables[j].Name = name
		variables[j].Type = VariableInteger
		variables[j].LowerBound = float64(bounds[j*2])
		variables[j].UpperBound = float64(bounds[j*2+1])
		variables[j].Index = j
		j++
	}
	return variables
}

/* Load a problem from a `.nl` file */
func ProblemFromFile(path string) (*Problem) {
	pathC := C.CString(path)
	asl := C.ASL_alloc(C.ASL_read_pfgh)
	nl := C.jac0dim_ASL(asl, pathC, C.ftnlen(len(path)))
	C.pfgh_read_ASL(asl, nl, C.ASL_find_o_class|C.ASL_find_c_class)
	aslPfgh := (*C.ASL_pfgh)(unsafe.Pointer(asl))
	return &Problem {path, asl, aslPfgh}
}

/* Return the larger integer */
func intMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
