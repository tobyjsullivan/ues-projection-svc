package projection

import (
    "github.com/tobyjsullivan/log-sdk/reader"
    "log"
    "os"
    "github.com/satori/go.uuid"
    "sync"
    "encoding/json"
)

const (
    EVENT_TYPE_ACCOUNT_OPENED = "AccountOpened"
)

var logger *log.Logger

func init() {
    logger = log.New(os.Stdout, "[projection] ", 0)
}

type Projection struct {
    mx sync.Mutex

    accounts map[uuid.UUID]*Account
}

func New() *Projection {
    return &Projection{
        accounts: make(map[uuid.UUID]*Account),
    }
}

func (p *Projection) Apply(e *reader.Event) {
    switch(e.Type) {
    case EVENT_TYPE_ACCOUNT_OPENED:
        p.handleAccountOpened(e.Data)
    }
}

func (p *Projection) FindAccount(accountId uuid.UUID) *Account {
    return p.accounts[accountId]
}

func (p *Projection) handleAccountOpened(data []byte) {
    var parsed struct {
        AccountID string `json:"accountId"`
    }
    if err := json.Unmarshal(data, &parsed); err != nil {
        logger.Println("Error processing AccountOpened event.", err.Error())
        return
    }

    var accountId uuid.UUID
    if err := accountId.UnmarshalText([]byte(parsed.AccountID)); err != nil {
        logger.Println("Error parsing account ID from AccountOpened event.", err.Error())
        return
    }

    acct := &Account{
        ID: accountId,
    }

    p.mx.Lock()
    defer p.mx.Unlock()

    p.accounts[acct.ID] = acct
}

