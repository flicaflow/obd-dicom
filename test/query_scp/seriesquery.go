package main

import (
	"database/sql"
	"git.onebytedata.com/odb/go-dicom/media"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type DCMSeries struct {
	SeriesDate        string
	SeriesTime        string
	Modality          string
	InstitutionName   string
	Description       string
	StudyInstanceUID  string
	SeriesInstanceUID string
	SeriesNumber      string
}

func (series *DCMSeries) Query(obj media.DcmObj) string {
	var tag media.DcmTag
	var query, partial string
	previous := false

	for i := 0; i < len(obj.GetTags()); i++ {
		tag = obj.GetTag(i)
		if tag.Length > 0 {
			switch tag.Group {
			case 0x08:
				switch tag.Element {
				case 0x21:
					series.SeriesDate = tag.GetString()
					partial = "SeriesDate='" + series.SeriesDate + "'"
					break
				case 0x32:
					series.SeriesTime = tag.GetString()
					partial = "SeriesTime='" + series.SeriesTime + "'"
					break
				case 0x60:
					series.Modality = tag.GetString()
					partial = "Modality='" + series.Modality + "'"
					break
				case 0x80:
					series.InstitutionName = tag.GetString()
					partial = "InstitutionName LIKE '%" + series.InstitutionName + "'%"
					break
				case 0x103E:
					series.Description = tag.GetString()
					partial = "Description LIKE '%" + series.Description + "%'"
					break
				}
				break
			case 0x20:
				switch tag.Element {
				case 0x0D:
					series.StudyInstanceUID = tag.GetString()
					partial = "StudyInstanceUID='" + series.StudyInstanceUID + "'"
					break
				case 0x0E:
					series.SeriesInstanceUID = tag.GetString()
					partial = "SeriesInstanceUID='" + series.SeriesInstanceUID + "'"
					break
				case 0x11:
					series.SeriesNumber = tag.GetString()
					partial = "SeriesNumber='" + series.StudyInstanceUID + "'"
					break
				}
				break
			}
			if len(partial) > 0 {
				if previous == true {
					query = query + " AND " + partial
				} else {
					query = " WHERE " + partial
					previous = true
				}
			}
		}
	}
	return query
}

func (series *DCMSeries) QueryResult(obj media.DcmObj) media.DcmObj {
	var added bool
	var tag media.DcmTag
	query := media.NewEmptyDCMObj()
	query.SetTransferSyntax("1.2.840.10008.1.2")

	for i := 0; i < len(obj.GetTags()); i++ {
		tag = obj.GetTag(i)
		added = true
		switch tag.Group {
		case 0x08:
			switch tag.Element {
			case 0x21:
				query.WriteStringGE(0x08, 0x21, "DA", series.SeriesDate)
				break
			case 0x32:
				query.WriteStringGE(0x08, 0x32, "TM", series.SeriesTime)
				break
			case 0x52:
				query.WriteStringGE(0x08, 0x52, "CS", "SERIES")
				break
			case 0x60:
				query.WriteStringGE(0x08, 0x60, "CS", series.Modality)
				break
			case 0x103E:
				query.WriteStringGE(0x08, 0x103E, "LO", series.Description)
				break
			default:
				added = false
			}
			break
		case 0x20:
			switch tag.Element {
			case 0x0D:
				query.WriteStringGE(0x20, 0x0d, "UI", series.StudyInstanceUID)
				break
			case 0x0E:
				query.WriteStringGE(0x20, 0x0e, "UI", series.SeriesInstanceUID)
				break
			case 0x11:
				query.WriteStringGE(0x20, 0x11, "IS", series.SeriesNumber)
				break
			default:
				added = false
			}
			break
		default:
			added = false
		}
		if added == false {
			query.Add(tag)
		}
	}
	return query
}

func (series *DCMSeries) Select(query media.DcmObj) (error, []media.DcmObj) {
	QueryString := series.Query(query)
	results := make([]media.DcmObj, 0)
	db, err := sql.Open("sqlite3", "./pacs.db")
	if err != nil {
		log.Println(err.Error())
		return err, nil
	}

	fields := "SeriesDate, SeriesTime, Modality, InstitutionName, SeriesDescription, StudyInstanceUID, SeriesInstanceUID, SeriesNumber"
	QueryString = "SELECT " + fields + " FROM Study " + QueryString
	rows, err := db.Query(QueryString)
	if err != nil {
		log.Println(err.Error())
		return err, nil
	}

	for rows.Next() {
		rows.Scan(&series.SeriesDate, &series.SeriesTime, &series.Modality, &series.InstitutionName, &series.Description, &series.StudyInstanceUID, &series.SeriesInstanceUID, &series.SeriesNumber)
		obj := series.QueryResult(query)
		results = append(results, obj)
	}
	rows.Close()
	return nil, results
}