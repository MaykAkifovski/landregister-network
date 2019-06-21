package main

import (
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
	Owners            []Owner           `json:"owners"`
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
type Owner struct {
	IdentityNumber string `json:"identityNumber"`
	Title          string `json:"title"` // Anrede
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	DateOfBirth    string `json:"dateOfBirth"`
	Postcode       string `json:"postcode"`
	City           string `json:"city"`
	Street         string `json:"street"`
	Streetnumber   string `json:"streetnumber"`
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
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "initLedger" { //create a new marble
		return t.initLedger(stub, args)
	} else if function == "queryLandRegister" {
		return t.queryLandRegister(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// initLedger - create a new realEstate, store into chaincode state
// ============================================================
func (t *SimpleChaincode) initLedger(stub shim.ChaincodeStubInterface, args []string) peer.Response {

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
			Owners: []Owner{
				{
					IdentityNumber: "1",
					Title:          "Mr",
					Firstname:      "Reiner",
					Lastname:       "Schatz",
					DateOfBirth:    "17.06.1955",
					Postcode:       "10***",
					City:           "Berlin",
					Street:         "Street",
					Streetnumber:   "123",
				},
				{
					IdentityNumber: "2",
					Title:          "Mrs",
					Firstname:      "Monika",
					Lastname:       "Schatz",
					DateOfBirth:    "16.07.1956",
					Postcode:       "10***",
					City:           "Berlin",
					Street:         "Street",
					Streetnumber:   "123",
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
	var grundbuchId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting grundbuchId to query")
	}

	grundbuchId = args[0]
	valAsBytes, err := stub.GetState(grundbuchId) //get the grundbuch from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + grundbuchId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsBytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + grundbuchId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsBytes)
}
