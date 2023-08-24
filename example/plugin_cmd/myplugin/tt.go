package myplugin

import (
	"time"

	"go.uber.org/zap"
)

type DenemeItem struct {
	Ben string
}

func (d *DenemeItem) SetStr(t string) {
	d.Ben = t
}

type TestItem struct {
	DenemeItem      *DenemeItem
	DenemeItemSlice []*DenemeItem
	Zaman           time.Time
}

func (t *TestItem) Init() error {
	return nil
}
func (t *TestItem) AddSlice(tn string) {
	dn := &DenemeItem{}
	dn.SetStr(tn)
	t.DenemeItemSlice = append(t.DenemeItemSlice, dn)
}
func (t *TestItem) GetSlice() []*DenemeItem {

	return t.DenemeItemSlice
}

func (t *TestItem) Name() string {
	return "Test"
}
func (t *TestItem) SetLogger(logger *zap.Logger) {

}
func (t *TestItem) CreateDn(tn string) {
	dn := &DenemeItem{}
	dn.SetStr(tn)
	t.DenemeItem = dn
}
func (t *TestItem) CreateDnX(tn string) string {
	dn := &DenemeItem{}
	dn.SetStr(tn)
	t.DenemeItem = dn
	return ""
}
func (t *TestItem) CreateDnXS(tn string) (string, string) {
	dn := &DenemeItem{}
	dn.SetStr(tn)
	t.DenemeItem = dn
	return "", ""
}
func (t *TestItem) CreateDnXSX(tn string) (string, string, error) {
	dn := &DenemeItem{}
	dn.SetStr(tn)
	t.DenemeItem = dn
	return "", "", nil
}
func (t *TestItem) CreateDnXSXX(tn string) (string, string, []string, error) {
	dn := &DenemeItem{}
	dn.SetStr(tn)
	t.DenemeItem = dn
	return "", "", nil, nil
}
