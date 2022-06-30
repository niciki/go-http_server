package main

type Statistic struct {
	DataEmployee   RecordEmployee
	DataDepartment []RecordDepartment
}

type RecordEmployee struct {
	CountOfEmployee int
	MedianAge       float64
	MedianSalary    float64
}

type RecordDepartment struct {
	DepartmentID int
	CountStaff   int
	MedianAge    float64
	MedianSalary float64
}

func NewStatistic(h *Handler) (Statistic, error) {
	EmployeeData, err := h.db.GetAll()
	if err != nil {
		return Statistic{}, err
	}
	AgeTotal := 0.0
	SalaryTotal := 0.0
	for _, j := range EmployeeData {
		AgeTotal += float64(j.Age)
		SalaryTotal += float64(j.Salary)
	}
	RecordEmpl := RecordEmployee{
		CountOfEmployee: len(EmployeeData),
		MedianAge:       AgeTotal / float64(len(EmployeeData)),
		MedianSalary:    SalaryTotal / float64(len(EmployeeData)),
	}
	DepartmentData, err := h.company.GetALL(&h.db)
	if err != nil {
		return Statistic{}, err
	}
	DataDep := make([]RecordDepartment, len(DepartmentData))
	for i, j := range DepartmentData {
		AgeTotal = 0.0
		SalaryTotal = 0.0
		for _, rec := range j.Staff {
			AgeTotal += float64(rec.Age)
			SalaryTotal += float64(rec.Salary)
		}
		DataDep[i] = RecordDepartment{
			DepartmentID: j.DepartmentID,
			CountStaff:   j.CountStaff,
			MedianAge:    AgeTotal / float64(j.CountStaff),
			MedianSalary: SalaryTotal / float64(j.CountStaff),
		}
	}
	return Statistic{DataEmployee: RecordEmpl, DataDepartment: DataDep}, nil
}
