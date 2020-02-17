package usecase

type IBranchUpdateParams interface {
	SetBranchName(name string) IBranchUpdateParams
	SetBranchAuthHost(name string) IBranchUpdateParams

	BranchName() (set bool, val string)
	BranchAuthHost() (set bool, val string)
}

type branchUpdateParams struct {
	branchName    string
	branchNameSet bool

	branchAuthHost    string
	branchAuthHostSet bool
}

func (b *branchUpdateParams) BranchAuthHost() (set bool, val string) {
	return b.branchAuthHostSet, b.branchAuthHost
}

func (b *branchUpdateParams) SetBranchAuthHost(branchAuthHost string) IBranchUpdateParams {
	b.branchAuthHostSet = true
	b.branchAuthHost = branchAuthHost

	return b
}

func (b *branchUpdateParams) BranchName() (set bool, val string) {
	return b.branchNameSet, b.branchName
}

func (b *branchUpdateParams) SetBranchName(branchName string) IBranchUpdateParams {
	b.branchNameSet = true
	b.branchName = branchName

	return b
}

func NewBranchUpdateParams() IBranchUpdateParams {
	return &branchUpdateParams{}
}
