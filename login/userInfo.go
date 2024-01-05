package login

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/cbrand/vodafone-billing-downloader/fetcher"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

const USER_INFO_URL = "https://api.vodafone.de/meinvodafone/v2/user/userInfo"

var (
	ErrUserInfoRequestFailed = errors.New("user info request failed")

	headerFmt = color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt = color.New(color.FgYellow).SprintfFunc()
)

func GetUserInfo(bearerToken fetcher.BearerToken) (*UserInfo, error) {
	var userInfo UserInfo
	err := fetcher.GetJson(USER_INFO_URL, bearerToken, &userInfo)
	if err == fetcher.ErrJsonRequestFailed {
		return nil, ErrUserInfoRequestFailed
	}
	return &userInfo, err
}

type UserInfo struct {
	UserAccountVBO *UserAccountVBO `json:"userAccountVBO"`
}

func (userInfo *UserInfo) HumanReadableString() string {
	return fmt.Sprintf("=== User Account ===\n%s", userInfo.UserAccountVBO.HumanReadableString())
}

func (userInfo *UserInfo) GetActiveContractCableID() string {
	return userInfo.UserAccountVBO.GetActiveContractCableID()
}

func (userInfo *UserInfo) GetAllContractIDs() []string {
	return userInfo.UserAccountVBO.GetAllContractIDs()
}

type UserAccountVBO struct {
	AuthLevel           string               `json:"authLevel"`
	OnlineUser          *OnlineUser          `json:"onlineUser,omitempty"`
	CableAccounts       []*CableAccount      `json:"cable,omitempty"`
	ActiveContractCable *ActiveContractCable `json:"activeContractCable,omitempty"`
}

func (userAccountVBO *UserAccountVBO) HumanReadableString() string {
	elements := []string{
		userAccountVBO.OnlineUser.HumanReadableString(),
		"=== Cable Accounts ===",
		userAccountVBO.HumanReadableContractTable(),
	}
	return strings.Join(elements, "\n")
}

func (UserAccountVBO *UserAccountVBO) HumanReadableContractTable() string {
	tbl := table.New("ID", "Name", "Active", "Number Subscriptions")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, cableAccount := range UserAccountVBO.CableAccounts {
		tbl.AddRow(cableAccount.ID, cableAccount.Name, cableAccount.IsActiveContract, len(cableAccount.Subscriptions))
	}

	data := bytes.NewBufferString("")
	tbl.WithWriter(data).Print()
	return data.String()
}

func (userAccountVBO *UserAccountVBO) GetActiveContractCableID() string {
	if userAccountVBO.ActiveContractCable == nil {
		return ""
	}
	return userAccountVBO.ActiveContractCable.IDString()
}

func (userAccountVBO *UserAccountVBO) GetAllContractIDs() []string {
	var contractIDs []string
	for _, cableAccount := range userAccountVBO.CableAccounts {
		contractIDs = append(contractIDs, cableAccount.ID)
	}
	return contractIDs
}

type OnlineUser struct {
	MintUserID            int    `json:"mintUserId"`
	UserName              string `json:"userName"`
	Title                 string `json:"title"`
	FirstName             string `json:"firstName"`
	LastName              string `json:"lastName"`
	LastLoginDate         string `json:"lastLoginDate"`
	PrimaryEmail          string `json:"primaryEmail"`
	EmailValidationStatus string `json:"emailValidationStatus"`
	IsFirstLogin          bool   `json:"isFirstLogin"`
	PermissionFlag        bool   `json:"permissionFlag"`
}

func (onlineUser *OnlineUser) HumanReadableString() string {
	tbl := table.New("Name", "Value")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	tbl.AddRow("MintUserID", onlineUser.MintUserID)
	tbl.AddRow("UserName", onlineUser.UserName)
	tbl.AddRow("Title", onlineUser.Title)
	tbl.AddRow("FirstName", onlineUser.FirstName)
	tbl.AddRow("LastName", onlineUser.LastName)
	tbl.AddRow("LastLoginDate", onlineUser.LastLoginDate)
	tbl.AddRow("PrimaryEmail", onlineUser.PrimaryEmail)
	tbl.AddRow("EmailValidationStatus", onlineUser.EmailValidationStatus)
	tbl.AddRow("IsFirstLogin", onlineUser.IsFirstLogin)
	tbl.AddRow("PermissionFlag", onlineUser.PermissionFlag)

	buffer := bytes.NewBufferString("")
	tbl.WithWriter(buffer).Print()
	return buffer.String()
}

type CableAccount struct {
	ID                string               `json:"id"`
	Name              string               `json:"name"`
	IsActiveContract  bool                 `json:"isActiveContract"`
	IsDefaultContract bool                 `json:"isDefaultContract"`
	HasCableMail      bool                 `json:"hasCableMail"`
	Subscriptions     []*CableSubscription `json:"subscription"`
}

type CableSubscription struct {
	ID            string `json:"id"`
	ActivatedDate string `json:"activatedDate"`
	Type          string `json:"type"`
	DisplayName   string `json:"displayName"`
}

type ActiveContractCable struct {
	ID   int    `json:"id"` // It is the same as `ID` in `CableSubscription` but as an integer because Vodafone is inconsistent
	Name string `json:"name"`
}

func (activeContractCable *ActiveContractCable) IDString() string {
	return fmt.Sprintf("%d", activeContractCable.ID)
}
