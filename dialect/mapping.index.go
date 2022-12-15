package dialect

import (
	"strings"
)

func (this_ *SqlMapping) GetIndexTypeInfos() (indexTypeInfoList []*IndexTypeInfo) {
	list := this_.indexTypeInfoList
	for _, one := range list {
		if one.IsExtend {
			continue
		}
		indexTypeInfoList = append(indexTypeInfoList, one)
	}
	return
}

func (this_ *SqlMapping) AddIndexTypeInfo(indexTypeInfo *IndexTypeInfo) {
	this_.indexTypeInfoCacheLock.Lock()
	defer this_.indexTypeInfoCacheLock.Unlock()

	if this_.indexTypeInfoCache == nil {
		this_.indexTypeInfoCache = make(map[string]*IndexTypeInfo)
	}

	key := strings.ToLower(indexTypeInfo.Name)
	find := this_.indexTypeInfoCache[key]
	this_.indexTypeInfoCache[key] = indexTypeInfo
	if find == nil {
		this_.indexTypeInfoList = append(this_.indexTypeInfoList, indexTypeInfo)
	} else {
		var list = this_.indexTypeInfoList
		var newList []*IndexTypeInfo
		for _, one := range list {
			if one == find {
				newList = append(newList, indexTypeInfo)
			} else {
				newList = append(newList, one)
			}
		}
		this_.indexTypeInfoList = newList
	}

	return
}

func (this_ *SqlMapping) GetIndexTypeInfo(typeName string) (indexTypeInfo *IndexTypeInfo, err error) {
	this_.indexTypeInfoCacheLock.Lock()
	defer this_.indexTypeInfoCacheLock.Unlock()

	if this_.indexTypeInfoCache == nil {
		this_.indexTypeInfoCache = make(map[string]*IndexTypeInfo)
	}

	key := strings.ToLower(typeName)
	indexTypeInfo = this_.indexTypeInfoCache[key]
	if indexTypeInfo == nil {

		indexTypeInfo = &IndexTypeInfo{
			Name:     typeName,
			Format:   typeName,
			IsExtend: true,
		}
		this_.indexTypeInfoCache[key] = indexTypeInfo
		//err = errors.New("dialect [" + this_.DialectType().Name + "] GetIndexTypeInfo not support index type name [" + typeName + "]")
		//fmt.Println(err)
		return
	}
	return
}

func (this_ *SqlMapping) IndexTypeFormat(index *IndexModel) (indexTypeFormat string, err error) {
	indexTypeInfo, err := this_.GetIndexTypeInfo(index.IndexType)
	if err != nil {
		return
	}
	if indexTypeInfo.IndexTypeFormat != nil {
		indexTypeFormat, err = indexTypeInfo.IndexTypeFormat(index)
		return
	}
	indexTypeFormat = indexTypeInfo.Format
	return
}

func (this_ *SqlMapping) IndexNameFormat(param *ParamModel, ownerName string, tableName string, index *IndexModel) (indexNameFormat string, err error) {
	if index.IndexName != "" {
		indexNameFormat = index.IndexName
		return
	}
	indexTypeInfo, err := this_.GetIndexTypeInfo(index.IndexType)
	if err != nil {
		return
	}
	if indexTypeInfo.IndexNameFormat != nil {
		indexNameFormat, err = indexTypeInfo.IndexNameFormat(param, ownerName, tableName, index)
		return
	}
	if ownerName != "" {
		indexNameFormat += ownerName + "_"
		//indexNameFormat += sortName(ownerName, 4) + "_"
	}
	if tableName != "" {
		indexNameFormat += tableName + "_"
		//indexNameFormat += sortName(tableName, 4) + "_"
	}
	if index.IndexType != "" && !strings.EqualFold(index.IndexType, "index") {
		indexNameFormat += index.IndexType + "_"
		//indexNameFormat += sortName(index.IndexType, 4) + "_"
	}
	indexNameFormat += strings.Join(index.ColumnNames, "_")
	//maxLength := 30 - len(indexNameFormat)
	//columnNamesStr := strings.Join(index.ColumnNames, "_")
	if this_.IndexNameMaxLen > 0 {
		indexNameFormat = sortName(indexNameFormat, this_.IndexNameMaxLen)
	}
	return
}

func sortName(name string, size int) (res string) {
	name = strings.TrimSpace(name)
	if len(name) <= size {
		res = name
		return
	}
	if strings.Contains(name, "_") {
		ss := strings.Split(name, "_")
		var names []string

		for _, s := range ss {
			if strings.TrimSpace(s) == "" {
				continue
			}
			names = append(names, strings.TrimSpace(s))
		}
		rSize := size / len(names)

		for i, s := range ss {
			if len(res) >= size {
				break
			}
			if i < len(ss)-1 {
				if rSize >= len(s)-1 {
					res += s + "_"
				} else {
					res += s[0:rSize-1] + "_"
				}
			} else {
				if rSize >= len(s) {
					res += s
				} else {
					res += s[0:rSize]
				}
			}
		}
	} else {
		res += name[0:size]
	}
	return
}
