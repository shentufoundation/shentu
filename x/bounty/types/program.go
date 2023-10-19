package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Programs is an array of program
type Programs []Program

func NewProgram(programId, name, adminAddr string, detail ProgramDetail, memberAddrs []string, status ProgramStatus) (Program, error) {

	return Program{
		ProgramId:      programId,
		Name:           name,
		Detail:         detail,
		AdminAddress:   adminAddr,
		MemberAccounts: memberAddrs,
		Status:         status,
	}, nil
}

func (m *Program) ValidateBasic() error {
	if len(m.ProgramId) == 0 {
		return ErrProgramID
	}
	if _, err := sdk.AccAddressFromBech32(m.AdminAddress); err != nil {
		return err
	}

	return nil
}

func NewDetail(desc, scopeRules, knownIssues string, bountyLevels []BountyLevel) ProgramDetail {

	return ProgramDetail{
		Description:  desc,
		ScopeRules:   scopeRules,
		KnownIssues:  knownIssues,
		BountyLevels: bountyLevels,
	}
}

//// UpdateDetail updates the fields of a given description. An error is
//// returned if the resulting description contains an invalid length.
//func (p ProgramDetail) UpdateDetail(p2 ProgramDetail) (ProgramDetail, error) {
//	if p2.Description  == DoNotModifyDesc {
//		d2.Moniker = d.Moniker
//	}
//
//	if d2.Identity == DoNotModifyDesc {
//		d2.Identity = d.Identity
//	}
//
//	if d2.Website == DoNotModifyDesc {
//		d2.Website = d.Website
//	}
//
//	if d2.SecurityContact == DoNotModifyDesc {
//		d2.SecurityContact = d.SecurityContact
//	}
//
//	if d2.Details == DoNotModifyDesc {
//		d2.Details = d.Details
//	}
//
//	return NewDescription(
//		d2.Moniker,
//		d2.Identity,
//		d2.Website,
//		d2.SecurityContact,
//		d2.Details,
//	).EnsureLength()
//}
