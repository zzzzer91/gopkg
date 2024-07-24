package langx

func FlattenMultiWayTree[T any](root T, getChildren func(T) []T) []T {
	var res []T
	q := []T{root}
	for len(q) > 0 {
		root = q[0]
		q = q[1:]
		res = append(res, root)
		children := getChildren(root)
		if len(children) > 0 {
			q = append(q, children...)
		}
	}
	return res
}
