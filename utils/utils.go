package utils

import "container/list"

type FIFOQueue struct {
    List *list.List
}

func NewFIFO() (fq *FIFOQueue) {
    return &FIFOQueue{
        List: list.New(),
    }
}

func (fq *FIFOQueue) Push(el string) {
    fq.List.PushBack(el)
}

func (fq *FIFOQueue) Pop() (v string) {
    el := fq.List.Front()
    if el != nil {
        v = el.Value.(string)
        fq.List.Remove(el)
    }
    return
}

func (fq *FIFOQueue) Len() (l int) {
    l = fq.List.Len()
    return
}
