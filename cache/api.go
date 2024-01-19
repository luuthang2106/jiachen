package cache

import (
	"jiachen/model"
	"jiachen/store"
)

var API = &Cache[string, struct{}]{
	// TTL: 1 * time.Second,
	// ReloadFunc: func(key string) (struct{}, error) {
	// 	api, err := store.API.FindOne(model.API{Path: key})
	// 	if err != nil {
	// 		return struct{}{}, err
	// 	}
	// 	if api.IsActive != nil && *api.IsActive {
	// 		return struct{}{}, nil
	// 	}
	// 	return struct{}{}, fmt.Errorf("API is not active")
	// },
}

func WarmupAPICache() (map[string]struct{}, error) {
	apis, err := store.API.Find(model.API{}, 0, 1000)
	if err != nil {
		return nil, err
	}
	m := make(map[string]struct{}, len(apis))
	for _, api := range apis {
		if api.IsActive != nil && *api.IsActive {
			m[api.Path] = struct{}{}
		}

	}
	return m, nil
}

// var API = (&Cache[string, struct{}]{}).
// 	// SetLoadFunc(
// 	// 	func() (map[string]struct{}, error) {
// 	// 		apis, err := store.API.Find(model.API{}, 0, 1000)
// 	// 		if err != nil {
// 	// 			return nil, err
// 	// 		}
// 	// 		m := make(map[string]struct{}, len(apis))
// 	// 		for _, api := range apis {
// 	// 			if api.IsActive != nil && *api.IsActive {
// 	// 				m[api.Path] = struct{}{}
// 	// 			}

// 	// 		}
// 	// 		return m, nil
// 	// 	},
// 	// ).
// 	SetLoadOneFunc(func(key string) (struct{}, error) {
// 		api, err := store.API.FindOne(model.API{Path: key})
// 		if err != nil {
// 			return struct{}{}, err
// 		}
// 		if api.IsActive != nil && *api.IsActive {
// 			return struct{}{}, nil
// 		}
// 		return struct{}{}, fmt.Errorf("API is not active")
// 	})
