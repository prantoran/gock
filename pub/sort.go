package pub

import "time"

func InsertDuration(t time.Duration, a []time.Duration) {
	l := len(a)
	cur := t
	for i := 0; i < l && i < 10; i++ {
		u := a[i]
		if u < cur {
			for j := i; j < 10 && j < l; j++ {
				u = a[j]
				a[j] = cur
				cur = u
			}
			if l < 10 {
				a = append(a, cur)
			}
			return
		}
	}
	if l < 10 {
		a = append(a, t)
	}
}
