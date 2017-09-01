package resources

import (
	"sync"
	"github.com/sarunask/awspec-generator/loggers"
)
//Tree of resources
type Tree struct {
	Tree map[string] *Resource
	Lock sync.Mutex
}


func (t *Tree) Init() {
	if t.Tree == nil {
		t.Lock.Lock()
		t.Tree = make(map[string]*Resource, 100)
		t.Lock.Unlock()
	}
}

func (t *Tree) Write(dir string, wg *sync.WaitGroup) {
	for i := range t.Tree {
		wg.Add(1)
		go func(res *Resource) {
			defer wg.Done()
			res.Write(dir)
		}(t.Tree[i])
	}
}

func (t *Tree) Push(r *Resource) {
	//Init tree if empty
	if t.Tree == nil {
		t.Init()
	}
	//Add resource by TerraformName to tree
	_, ok := t.Tree[r.TerraformName]
	if ok == false {
		if t.Tree == nil {
			loggers.Error.Println("Tree is still nil")
		}
		t.Lock.Lock()
		t.Tree[r.TerraformName] = r
		t.Lock.Unlock()
	}
}
