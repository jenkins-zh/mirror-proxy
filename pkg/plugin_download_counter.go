package pkg

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"time"
)

// GitPluginDownloadCounter count the data by git
type GitPluginDownloadCounter struct {
	Path     string
	Username string
	Password string
}

// PluginDownloadCounter
type PluginDownloadCounter interface {
	ReadData() ([]PluginDownloadData, error)
	FindByYear(year string) (*PluginDownloadData, error)
	Save(data *PluginDownloadData) error

	UpdateCenterCountIncrease(downloadData *PluginDownloadData) error
}

// ReadData get all data
func (g *GitPluginDownloadCounter) ReadData() (dataArray []PluginDownloadData, err error) {
	return
}

// FindByYear returns the data by year
func (g *GitPluginDownloadCounter) FindByYear(year string) (downloadData *PluginDownloadData, err error) {
	dataFilePath := g.GetDataFilePath(year)
	downloadData = &PluginDownloadData{}

	if _, err = os.Stat(dataFilePath); err != nil {
		return
	}

	var data []byte
	if data, err = ioutil.ReadFile(dataFilePath); err == nil {
		err = yaml.Unmarshal(data, downloadData)
	}
	return
}

// Save stores the data
func (g *GitPluginDownloadCounter) Save(downloadData *PluginDownloadData) (err error) {
	dataFilePath := g.GetDataFilePath(downloadData.Year)

	fmt.Println("prepare to store", dataFilePath)

	if err = os.MkdirAll(path.Dir(dataFilePath), 0751); err != nil {
		return
	}

	var data []byte
	if data, err = yaml.Marshal(downloadData); err == nil {
		err = ioutil.WriteFile(dataFilePath, data, 0644)
	}
	return
}

func (g *GitPluginDownloadCounter) UpdateCenterCountIncrease(downloadData *PluginDownloadData) (err error) {
	err = g.PluginCountIncrease(downloadData, "update-center")
	return
}

func (g *GitPluginDownloadCounter) PluginCountIncrease(downloadData *PluginDownloadData, plugin string) (err error) {
	if center, ok := downloadData.Plugins[plugin]; ok {
		if count, ok := center.Data[GetDate()]; ok {
			center.Data[GetDate()] = count + 1
		} else {
			center.Data[GetDate()] = 1
		}
	} else {
		if len(downloadData.Plugins) == 0 {
			downloadData.Plugins = make(map[string]PluginData, 1)
		}

		downloadData.Plugins[plugin] = PluginData{
			Data: map[string]int64{
				GetDate(): 1,
			},
		}
	}
	return
}

func (g *GitPluginDownloadCounter) RecordPluginDownloadData(plugin, provider string) (err error) {
	fmt.Println("plugin", plugin, "provider", provider)
	var downloadData *PluginDownloadData
	if downloadData, err = g.FindByYear(GetCurrentYear()); err != nil {
		fmt.Println("cannot find by year", GetCurrentYear(), err)
		downloadData = &PluginDownloadData{
			Year: GetCurrentYear(),
		}
	}
	fmt.Println(downloadData)

	if err = g.PluginCountIncrease(downloadData, plugin); err != nil {
		fmt.Println(err)
	}

	fmt.Println("pluginDownloadCounter.Save", downloadData)
	if err = g.Save(downloadData); err != nil {
		fmt.Println(err)
	}
	return
}

func (g *GitPluginDownloadCounter) RecordUpdateCenterVisitData() (err error) {
	var downloadData *PluginDownloadData
	if downloadData, err = g.FindByYear(GetCurrentYear()); err != nil {
		fmt.Println("cannot find by year", GetCurrentYear(), err)
		downloadData = &PluginDownloadData{
			Year: GetCurrentYear(),
		}
	}

	if err = g.UpdateCenterCountIncrease(downloadData); err != nil {
		fmt.Println(err)
	}

	fmt.Println("pluginDownloadCounter.Save", downloadData)
	if err = g.Save(downloadData); err != nil {
		fmt.Println(err)
	}
	return
}

func GetCurrentYear() string {
	dt := time.Now()
	return dt.Format("2006")
}

func GetDate() string {
	dt := time.Now()
	return dt.Format("2006-01-02")
}

func (g *GitPluginDownloadCounter) GetDataFilePath(year string) string {
	return path.Join(g.Path, fmt.Sprintf("%s.yaml", year))
}
