// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

import (
	"strconv"
)

// Account provides access to the account endpoints of the API.
//
// Use this to read the account, its balance and its usage, and to pause or resume
// its scheduled checks.
type Account struct {
	handler *RequestHandler
}

// Get retrieves the account: its balance, whether its scheduled checks are running,
// and how much of what the plan allows is in use.
//
// Example:
//
//	account, err := client.Account().Get()
func (a *Account) Get(opts ...RequestOption) (*Response, error) {
	return a.handler.Get("account", nil, opts...)
}

// Balance retrieves the account balance.
//
// Example:
//
//	balance, err := client.Account().Balance()
func (a *Account) Balance(opts ...RequestOption) (float64, error) {
	body, err := a.handler.GetPlain("account/balance", opts...)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(body, 64)
}

// Usage retrieves the account's usage by month.
//
// Optional parameters:
//   - months: How many months to return, the current month included
//
// Example:
//
//	usage, err := client.Account().Usage(map[string]interface{}{"months": 3})
func (a *Account) Usage(params map[string]interface{}, opts ...RequestOption) (*Response, error) {
	return a.handler.Get("account/usage", params, opts...)
}

// PauseChecks pauses the account's scheduled checks.
//
// Example:
//
//	result, err := client.Account().PauseChecks()
func (a *Account) PauseChecks(opts ...RequestOption) (*Response, error) {
	return a.handler.Post("account/checks/pause", nil, opts...)
}

// ResumeChecks resumes the account's scheduled checks.
//
// Example:
//
//	result, err := client.Account().ResumeChecks()
func (a *Account) ResumeChecks(opts ...RequestOption) (*Response, error) {
	return a.handler.Post("account/checks/resume", nil, opts...)
}
