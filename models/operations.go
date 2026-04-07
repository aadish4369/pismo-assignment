package models

type OperationType int

const (
	NormalPurchase      OperationType = 1
	InstallmentPurchase OperationType = 2
	Withdraw            OperationType = 3
	CreditVoucher       OperationType = 4
	InstallmentEMI      OperationType = 5
)
