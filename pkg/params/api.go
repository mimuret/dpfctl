package params

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis"
)

var (
	groupMap = NewGroupMap()
)

type GroupMap struct {
	groupMap map[string]APISetSlice
	cmdMap   map[string]*APISet
}

func NewGroupMap() *GroupMap {
	return &GroupMap{
		groupMap: make(map[string]APISetSlice),
		cmdMap:   make(map[string]*APISet),
	}
}

func (g *GroupMap) IterateAPISet(f func(apiSet *APISet), action api.Action) {
	for _, list := range g.groupMap {
		if list != nil {
			list.IterateAPISet(f, action)
		}
	}
}

func UpdateCmdMap() { groupMap.UpdateCmdMap() }

func (g *GroupMap) UpdateCmdMap() {
	newMap := make(map[string]*APISet)
	g.IterateAPISet(func(apiSet *APISet) {
		if _, ok := newMap[apiSet.Name]; ok {
			panic(fmt.Sprintf("cmd %s is duplicated", apiSet.Name))
		}
		newMap[apiSet.Name] = apiSet
	}, "")
	g.cmdMap = newMap
}
func ValidateParams() { groupMap.ValidateParams() }
func (g *GroupMap) ValidateParams() {
	for _, apiSet := range g.cmdMap {
		if apiSet == nil {
			continue
		}
		for action, apiSepc := range apiSet.Action {
			if apiSepc.Params != nil {
				if err := apiSepc.Params.Validate(); err != nil {
					panic(fmt.Sprintf("%s %s %s", apiSet.Name, action, err.Error()))
				}
			}
		}
	}
}

func GetAPISetfromCmdName(name string) *APISet { return groupMap.GetAPISetfromCmdName(name) }

func (g *GroupMap) GetAPISetfromCmdName(name string) *APISet { return g.cmdMap[name] }

func SetGroup(name string, list APISetSlice) {
	groupMap.SetGroup(name, list)
}

func (g *GroupMap) SetGroup(name string, list APISetSlice) {
	g.groupMap[name] = list
	g.UpdateCmdMap()
}

func GetGroupList() []string {
	return groupMap.GetGroupList()
}

func (g *GroupMap) GetGroupList() []string {
	groups := []string{}
	for groupName := range g.groupMap {
		groups = append(groups, groupName)
	}
	return groups
}

func GetAPISlice(name string) APISetSlice {
	return groupMap.GetAPISlice(name)
}

func (g *GroupMap) GetAPISlice(name string) APISetSlice {
	return g.groupMap[name]
}

func GetValidArgs(action api.Action) []string {
	return groupMap.GetValidArgs(action)
}

func (g *GroupMap) GetValidArgs(action api.Action) []string {
	res := []string{}
	for _, apiSet := range g.cmdMap {
		if apiSet != nil {
			if _, ok := apiSet.Action[action]; ok {
				res = append(res, apiSet.Name)
			}
		}
	}
	return res
}

type APISetSlice []*APISet

func (a *APISetSlice) IterateAPISet(f func(apiSet *APISet), action api.Action) {
	for _, apiSet := range *a {
		if apiSet != nil {
			_, ok := apiSet.Action[action]
			if action == "" || ok {
				f(apiSet)
			}
		}
	}
}

type APISet struct {
	Name        string
	Description string
	Action      map[api.Action]API
}

type API struct {
	Object api.Spec
	Params Params
	// set args for apply
	SetFunc func(apis.Spec, []string) error
}

type Params []Param

func (p Params) String() string {
	paramString := []string{}
	for _, param := range p {
		paramString = append(paramString, param.String())
	}
	return strings.Join(paramString, " ")
}

func (p Params) Validate() error {
	for i, param := range p {
		if err := param.Validate(); err != nil {
			return fmt.Errorf("param[%d]: %w", i, err)
		}
		if param.Type == ParamTypeArrayString && i != len(p)-1 {
			return fmt.Errorf("param[%d]: ParamTypeArrayString must be last item", i)
		}
		if !param.Required && i != len(p)-1 {
			return fmt.Errorf("param[%d]: Required false must be last item", i)
		}
	}
	return nil
}

func (p Params) MakeArgs(args []string) ([]interface{}, error) {
	sets := []interface{}{}
	i := 0
	j := 0
	enough := false
	for i = 0; i < len(args) && j < len(p); i++ {
		v, err := p[j].Parse(args[i])
		if err != nil {
			return nil, fmt.Errorf("param[%d]: %w", i, err)
		}
		sets = append(sets, v)
		if p[j].Type == ParamTypeArrayString || p[j].Type == ParamTypeArrayInt64 {
			enough = true
		} else {
			j++
		}
	}
	if !enough {
		if j < len(p) && p[j].Required {
			return nil, fmt.Errorf("not enough parameters, %d %d", j, len(p))
		}
	}
	return sets, nil
}

func (p Params) SetArgs(spec apis.Spec, args []string) error {
	sets, err := p.MakeArgs(args)
	if err != nil {
		return err
	}
	return spec.SetParams(sets...)
}

type ParamType int

const (
	ParamTypeInt64       ParamType = 1
	ParamTypeArrayInt64  ParamType = 2
	ParamTypeString      ParamType = 3
	ParamTypeArrayString ParamType = 4
)

type Param struct {
	Name     string
	Required bool
	Type     ParamType
}

func (p Param) Validate() error {
	if p.Name == "" {
		return errors.New("name must not be empty")
	}
	if strings.ContainsAny(p.Name, "[]./ $%^&*()!{}|\"':?><") {
		return errors.New("name must not include special character")
	}
	switch p.Type {
	case ParamTypeInt64, ParamTypeArrayInt64, ParamTypeString, ParamTypeArrayString:
	default:
		return errors.New("type error")
	}
	return nil
}

func (p Param) Parse(s string) (interface{}, error) {
	switch p.Type {
	case ParamTypeInt64, ParamTypeArrayInt64:
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("int64 parse error")
		}
		return v, nil
	case ParamTypeString, ParamTypeArrayString:
		return s, nil
	}
	return nil, errors.New("type error")
}

func (p Param) String() string {
	res := p.Name
	if p.Type == ParamTypeArrayString {
		res += " ..."
	}
	if !p.Required {
		if p.Type == ParamTypeArrayString {
			res = "[" + res + "]"
		}
	}
	return res
}
