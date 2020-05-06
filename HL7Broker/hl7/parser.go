package hl7

// HL7 parser class
import (
	"rafael/DICOM/media"
	"strings"
)

var (
	study  media.DCMStudy
	params [64]string
)

// ParsePID - parses the PID Segment
func ParsePID(line string) {
	params := strings.Split(line, "|")
	study.PatientID = params[3]
	index := strings.Index(study.PatientID, "^")
	if index > 0 {
		study.PatientID = study.PatientID[:index]
	}
	study.PatientName = params[5]
	study.PatientBD = params[7]
	study.PatientSex = params[8]
}

// ParsePV1 - parses the PV1 Segment
func ParsePV1(line string) {
	params := strings.Split(line, "|")
	study.ReferringPhysician = params[8]
	index := strings.Index(study.ReferringPhysician, "^")
	if index > 0 {
		study.ReferringPhysician = study.ReferringPhysician[index:]
	}
}

// ParseORC - parses the ORC Segment
func ParseORC(line string) {
	params := strings.Split(line, "|")
	tDate := strings.Split(params[7], "^")
	if tDate[4] != "" {
		study.ReportTime = tDate[4]
		if len(study.ReportTime) > 14 {
			study.ReportDate = study.ReportTime[0:7]
			study.ReportTime = study.ReportTime[8:14]
		}
	}
}

// ParseOBR - parses the OBR Segment
func ParseOBR(line string) {
	params := strings.Split(line, "|")
	if study.ReportDate == "" {
		study.ReportTime = params[6]
		if len(study.ReportTime) > 14 {
			study.ReportDate = study.ReportTime[0:7]
			study.ReportTime = study.ReportTime[8:14]
		}
	}
	study.AccessionNumber = params[18]
	study.Modality = params[24]
	study.Description = params[44]
	if study.Description == "" {
		study.Description = params[4]
	}
}

// ParseOBX - parses the OBX Segment
func ParseOBX(line string) {
	params := strings.Split(line, "|")
	if params[2] == "TX" {
		study.ReportText = study.ReportText + params[5]
	}
}

// SaveDICOMSR - Save ORU to DICOM format
func SaveDICOMSR(fileName string) {
	var srobj media.DcmObj

	srobj.ExplicitVR = true
	srobj.BigEndian = false
	srobj.TransferSyntax = "1.2.840.10008.1.2.1"
	study.StudyInstanceUID = "9999.9999.1"
	srobj.CreateSR(study, "8888.8888.1", "7777.7777.1")
	srobj.Write(fileName)
}
