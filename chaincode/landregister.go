package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Grundbuch
type landRegister struct {
	ObjectType        string            `json:"docType"` //docType is used to distinguish the various types of objects in state database
	TitlePage         titlePage         `json:"titlePage"`
	InventoryRegister inventoryRegister `json:"inventoryRegister"`
	Owners            []owner           `json:"owners"`
	ReservationNote   bool              `json:"reservationNote"`
}

// Titelblatt
type titlePage struct {
	DistrictCourt        string `json:"districtCourt"`        // Amtsgericht
	LandRegistryDistrict string `json:"landRegistryDistrict"` // Grundbuchbezirk
	SheetNumber          string `json:"sheetNumber"`          // Nummer des Blattes
}

// Bestandsverzeichnis
type inventoryRegister struct {
	Subdistrict  string `json:"subdistrict"`  // Gemarkung
	Hall         string `json:"hall"`         // Flur
	Parcel       string `json:"parcel"`       // Flurstueck
	EconomicType string `json:"economicType"` // Wirtschaftsart
	Location     string `json:"location"`     // Lage
	Size         string `json:"size"`         // Groesse
}

// Eigentuemer
type owner struct {
	IdentityNumber string `json:"identityNumber"`
	Title          string `json:"title"` // Anrede
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	DateOfBirth    string `json:"dateOfBirth"`
	Postcode       string `json:"postcode"`
	City           string `json:"city"`
	Street         string `json:"street"`
	StreetNumber   string `json:"streetNumber"`
}

type reservationNoteRequest struct {
	ObjectType        string            `json:"docType"` //docType is used to distinguish the various types of objects in state database
	TitlePage         titlePage         `json:"titlePage"`
	InventoryRegister inventoryRegister `json:"inventoryRegister"`
	Owners            []owner           `json:"owners"`
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return t.initLedger(stub)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "initLedger" { //create a new marble
		return t.initLedger(stub)
	} else if function == "queryLandRegister" {
		return t.queryLandRegister(stub, args)
	} else if function == "queryAllLandRegisters" {
		return t.queryAllLandRegisters(stub, args)
	} else if function == "createLandRegister" {
		return t.createLandRegister(stub, args)
	} else if function == "createReservationNote" {
		return t.createReservationNote(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// initLedger - create a new realEstate, store into chaincode state
// ============================================================
func (t *SimpleChaincode) initLedger(stub shim.ChaincodeStubInterface) peer.Response {

	landRegisters := []landRegister{
		{
			ObjectType: "landRegister",
			TitlePage: titlePage{
				DistrictCourt:        "Eutin",
				LandRegistryDistrict: "Malente",
				SheetNumber:          "3323",
			},
			InventoryRegister: inventoryRegister{
				Subdistrict:  "Malente",
				Hall:         "4",
				Parcel:       "6/12",
				EconomicType: "Hof- und Gebaeudeflaeche",
				Location:     "Steencamp 112",
				Size:         "845 m2",
			},
			Owners: []owner{
				{
					IdentityNumber: "1",
					Title:          "Mr",
					FirstName:      "Reiner",
					LastName:       "Schatz",
					DateOfBirth:    "17.06.1955",
					Postcode:       "10***",
					City:           "Berlin",
					Street:         "Street",
					StreetNumber:   "123",
				},
				{
					IdentityNumber: "2",
					Title:          "Mrs",
					FirstName:      "Monika",
					LastName:       "Schatz",
					DateOfBirth:    "16.07.1956",
					Postcode:       "10***",
					City:           "Berlin",
					Street:         "Street",
					StreetNumber:   "123",
				},
			},
			ReservationNote: false,
		},
	}

	i := 0
	for i < len(landRegisters) {
		fmt.Println("i is ", i)
		landRegisterAsBytes, err := json.Marshal(landRegisters[i])
		if err != nil {
			return shim.Error(err.Error())
		}
		landRegisterId := landRegisters[i].ObjectType + "-" +
			landRegisters[i].TitlePage.DistrictCourt + "-" +
			landRegisters[i].TitlePage.LandRegistryDistrict + "-" +
			landRegisters[i].TitlePage.SheetNumber
		err = stub.PutState(landRegisterId, landRegisterAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		/*
			CompositeIndex
		*/

		fmt.Println("Added ", landRegisters[i], " with id: ", landRegisterId)
		i = i + 1
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) queryLandRegister(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var landRegisterId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting landRegisterId to query")
	}

	landRegisterId = args[0]
	valAsBytes, err := stub.GetState(landRegisterId) //get the grundbuch from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + landRegisterId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsBytes == nil {
		jsonResp = "{\"Error\":\"Landregister does not exist: " + landRegisterId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsBytes)
}

func (t *SimpleChaincode) queryAllLandRegisters(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var jsonResp string
	var err error

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0 arguments.")
	}

	resultsIterator, err := stub.GetStateByRange("", "") // get the all Landregisters
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get all states\"}"
		return shim.Error(jsonResp)
	} else if resultsIterator == nil {
		jsonResp = "{\"Error\":\"Landregisters do not exist\"}"
		return shim.Error(jsonResp)
	} else {
		defer resultsIterator.Close()

		// buffer is a JSON array containing QueryResults
		var buffer bytes.Buffer
		buffer.WriteString("[")

		bArrayMemberAlreadyWritten := false
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return shim.Error(err.Error())
			}
			// Add a comma before array members, suppress it for the first array member
			if bArrayMemberAlreadyWritten == true {
				buffer.WriteString(",")
			}
			buffer.WriteString("{\"Key\":")
			buffer.WriteString("\"")
			buffer.WriteString(queryResponse.Key)
			buffer.WriteString("\"")

			buffer.WriteString(", \"Record\":")
			// Record is a JSON object, so we write as-is
			buffer.WriteString(string(queryResponse.Value))
			buffer.WriteString("}")
			bArrayMemberAlreadyWritten = true
		}
		buffer.WriteString("]")

		fmt.Printf("- queryAllLandRegisters:\n%s\n", buffer.String())

		return shim.Success(buffer.Bytes())
	}
}

func (t *SimpleChaincode) createLandRegister(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var landRegisterAsString, jsonResp string
	var landRegisterAsObject landRegister
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting JSON Representation of LandRegister to query")
	}

	landRegisterAsString = args[0]
	err = json.Unmarshal([]byte(landRegisterAsString), &landRegisterAsObject)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to unmarshal: " + landRegisterAsString + "\"}"
		return shim.Error(jsonResp)
	}

	landRegisterId := landRegisterAsObject.ObjectType + "-" +
		landRegisterAsObject.TitlePage.DistrictCourt + "-" +
		landRegisterAsObject.TitlePage.LandRegistryDistrict + "-" +
		landRegisterAsObject.TitlePage.SheetNumber
	err = stub.PutState(landRegisterId, []byte(landRegisterAsString))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Added ", landRegisterAsString, " with id: ", landRegisterId)

	return shim.Success(nil)
}

func (t *SimpleChaincode) createReservationNote(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var reservationNoteRequestAsString, jsonResp string
	var reservationNoteRequestAsObject reservationNoteRequest
	var landRegisterAsObject landRegister
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting JSON Representation of ReservationNoteRequest to query")
	}

	reservationNoteRequestAsString = args[0]
	err = json.Unmarshal([]byte(reservationNoteRequestAsString), &reservationNoteRequestAsObject)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to unmarshal: " + reservationNoteRequestAsString + "\"}"
		return shim.Error(jsonResp)
	}

	landRegisterId := reservationNoteRequestAsObject.ObjectType + "-" +
		reservationNoteRequestAsObject.TitlePage.DistrictCourt + "-" +
		reservationNoteRequestAsObject.TitlePage.LandRegistryDistrict + "-" +
		reservationNoteRequestAsObject.TitlePage.SheetNumber

	landRegisterAsBytes, err := stub.GetState(landRegisterId)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to find Landregister for: " + landRegisterId + "\"}"
		return shim.Error(jsonResp)
	}
	err = json.Unmarshal(landRegisterAsBytes, &landRegisterAsObject)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to unmarshal: " + string(landRegisterAsBytes) + "\"}"
		return shim.Error(jsonResp)
	}

	if landRegisterAsObject.ReservationNote {
		jsonResp = "{\"Error\":\"Reservation Note already exists!\"}"
		return shim.Error(jsonResp);
	}

	if !assertSameLandRegister(landRegisterAsObject, reservationNoteRequestAsObject) {
		jsonResp = "{\"Error\":\"Landregister and Reservation note request unequal!\"}"
		return shim.Error(jsonResp)
	}
	landRegisterAsObject.ReservationNote = true

	landRegisterAsBytes, _ = json.Marshal(landRegisterAsObject)
	err = stub.PutState(landRegisterId, landRegisterAsBytes)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to create ReservationNote for " + landRegisterId + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(nil)
}

func assertSameLandRegister(landRegister landRegister, request reservationNoteRequest) bool {
	if assertSameInventoryRegister(landRegister, request) {
	} else {
		return false
	}

	if len(landRegister.Owners) != len(request.Owners) {
		return false
	}

	numberOfOwners := len(landRegister.Owners)
	for i := 0; i < numberOfOwners; i++ {
		for j := 0; j < numberOfOwners; j++ {
			if landRegister.Owners[i].IdentityNumber == request.Owners[j].IdentityNumber {
				sameOwner := assertSameOwner(landRegister.Owners[i], request.Owners[j])
				if !sameOwner {
					return false
				}
			}
		}
	}

	return true
}

func assertSameInventoryRegister(landRegister landRegister, request reservationNoteRequest) bool {
	return landRegister.InventoryRegister.Subdistrict == request.InventoryRegister.Subdistrict &&
		landRegister.InventoryRegister.Size == request.InventoryRegister.Size &&
		landRegister.InventoryRegister.Location == request.InventoryRegister.Location &&
		landRegister.InventoryRegister.EconomicType == request.InventoryRegister.EconomicType &&
		landRegister.InventoryRegister.Parcel == request.InventoryRegister.Parcel &&
		landRegister.InventoryRegister.Hall == request.InventoryRegister.Hall
}

func assertSameOwner(lrOwner owner, rOwner owner) bool {
	return lrOwner.DateOfBirth == rOwner.DateOfBirth &&
		lrOwner.City == rOwner.City &&
		lrOwner.Postcode == rOwner.Postcode &&
		lrOwner.StreetNumber == rOwner.StreetNumber &&
		lrOwner.Title == rOwner.Title &&
		lrOwner.LastName == rOwner.LastName &&
		lrOwner.FirstName == rOwner.FirstName
}
