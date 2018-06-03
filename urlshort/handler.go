package urlshort

import (
	"net/http"
	yaml "gopkg.in/yaml.v2"
)

type pathUrl struct {
	Path string `yaml:"path"`
	URL string	`yaml:"URL"`
}

//MapHandler will return an http.Handler func that will attempt to map any paths
//in the map to their corresponding URL values. If the path is not provided in the map,
//then the http.Handler will be called instead.
func MapHandler(pathsToURLs map[string]string, fallback http.Handler) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		path:=r.URL.Path
		if dest, ok:= pathsToURLs[path]; ok{
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		//if we can math a path---
		//	redirect to it
		// 		otherwise use the fallback
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error){
	//parse the yaml
	pathUrls, err:= parseYaml(yamlBytes)
	if err != nil {
		return nil, err
		}
		pathsToUrls:= buildMap(pathUrls)
		return MapHandler(pathsToUrls, fallback), nil
		}

func buildMap(pathUrls []pathUrl) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu:= range pathUrls{
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}

func parseYaml(data []byte) ([]pathUrl, error){
	var pathUrls []pathUrl
	err:= yaml.Unmarshal(data, &pathUrls)
	if err!=nil{
		return nil, err
	}
	return pathUrls, nil
}

