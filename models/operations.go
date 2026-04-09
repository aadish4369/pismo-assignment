package models

type OperationType int

const (
	NormalPurchase      OperationType = 1
	InstallmentPurchase OperationType = 2
	Withdraw            OperationType = 3
	CreditVoucher       OperationType = 4
)

func (operationType OperationType) IsValid() bool {
	switch operationType {
	case NormalPurchase, InstallmentPurchase, Withdraw, CreditVoucher:
		return true
	default:
		return false
	}
}

func (operationType OperationType) IsCredit() bool {
	return operationType == CreditVoucher
}
