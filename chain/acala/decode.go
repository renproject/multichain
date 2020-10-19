package acala

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

type eventMinted struct {
	Phase    types.Phase
	Who      types.AccountID
	Currency types.U8
	Amount   types.U128
	Topics   []types.Hash
}

type eventBurnt struct {
	Owner  types.AccountID
	Dest   types.AccountID
	Amount types.U128
}

type EventsWithMint struct {
	types.EventRecords
	RenToken_Minted []eventMinted //nolint:stylecheck,golint
}

type EventsWithBurn struct {
	types.EventRecords
	RenToken_Burnt []eventBurnt //nolint:stylecheck,golint
}

// ParseEvents decodes the events records from an EventRecordRaw into a target
// t using the given Metadata m. If this method returns an error like `unable
// to decode Phase for event #x: EOF`, it is likely that you have defined a
// custom event record with a wrong type. For example your custom event record
// has a field with a length prefixed type, such as types.Bytes, where your
// event in reality contains a fixed width type, such as a types.U32.
func ParseEvents(e *types.EventRecordsRaw, m *types.Metadata, t interface{}) error {
	// ensure t is a pointer
	ttyp := reflect.TypeOf(t)
	if ttyp.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer, but is " + fmt.Sprint(ttyp))
	}
	// ensure t is not a nil pointer
	tval := reflect.ValueOf(t)
	if tval.IsNil() {
		return errors.New("target is a nil pointer")
	}
	val := tval.Elem()
	typ := val.Type()
	// ensure val can be set
	if !val.CanSet() {
		return fmt.Errorf("unsettable value %v", typ)
	}
	// ensure val points to a struct
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("target must point to a struct, but is " + fmt.Sprint(typ))
	}

	decoder := scale.NewDecoder(bytes.NewReader(*e))

	// determine number of events
	n, err := decoder.DecodeUintCompact()
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("found %v events", n))

	// iterate over events
	for i := uint64(0); i < n.Uint64(); i++ {
		fmt.Println(fmt.Sprintf("decoding event #%v", i))

		// decode Phase
		phase := types.Phase{}
		err := decoder.Decode(&phase)
		if err != nil {
			return fmt.Errorf("unable to decode Phase for event #%v: %v", i, err)
		}

		// decode EventID
		id := types.EventID{}
		err = decoder.Decode(&id)
		if err != nil {
			return fmt.Errorf("unable to decode EventID for event #%v: %v", i, err)
		}

		fmt.Println(fmt.Sprintf("event #%v has EventID %v", i, id))

		// ask metadata for method & event name for event
		moduleName, eventName, err := m.FindEventNamesForEventID(id)
		// moduleName, eventName, err := "System", "ExtrinsicSuccess", nil
		if err != nil {
			fmt.Printf("unable to find event with EventID %v in metadata for event #%v: %s\n", id, i, err)
			continue
			// return fmt.Errorf("unable to find event with EventID %v in metadata for event #%v: %s", id, i, err)
		}

		fmt.Println(fmt.Sprintf("event #%v is in module %v with event name %v", i, moduleName, eventName))

		// check whether name for eventID exists in t
		field := val.FieldByName(fmt.Sprintf("%v_%v", moduleName, eventName))
		if !field.IsValid() {
			eventParams, err := findEventForEventID(m.AsMetadataV10, id)
			if err != nil {
				return fmt.Errorf("unable to find event with EventID %v in metadata for event #%v: %s", id, i, err)
			}

			for j := 0; j < len(eventParams.Args); j++ {
				fmt.Printf("decoding field: %v (%v)\n", j, eventParams.Args[j])
				switch eventParams.Args[j] {
				case "u8":
					param := types.U8(0)
					err = decoder.Decode(param)
				case "u16":
					param := types.U16(0)
					err = decoder.Decode(param)
				case "u32":
					param := types.U32(0)
					err = decoder.Decode(param)
				case "u64":
					param := types.U64(0)
					err = decoder.Decode(param)
				case "u128":
					param := types.U128{}
					err = decoder.Decode(param)
				case "u256":
					param := types.U256{}
					err = decoder.Decode(param)
				case "Phase":
					param := types.Phase{}
					err = decoder.Decode(param)
				case "DispatchInfo":
					param := types.DispatchInfo{}
					err = decoder.Decode(param)
				case "DispatchError":
					param := types.DispatchError{}
					err = decoder.Decode(param)
				case "AccountId":
					param := types.AccountID{}
					err = decoder.Decode(param)
				case "AccountIndex":
					param := types.AccountIndex(0)
					err = decoder.Decode(param)
				// case "Balance":
				// 	param := types.Balance{}
				// 	err = decoder.Decode(param)
				// case "Status":
				// 	param := types.Status{}
				// 	err = decoder.Decode(param)
				case "bool":
					param := types.Bool(false)
					err = decoder.Decode(param)
				// case "CallHash":
				// 	param := types.CallHash{}
				// 	err = decoder.Decode(param)
				// case "Timepoint":
				// param := types.Timepoint<BlockNumber>{}
				// err = decoder.Decode(param)
				// case "ProposalIndex":
				// 	param := types.ProposalIndex{}
				// 	err = decoder.Decode(param)
				case "Hash":
					param := types.Hash{}
					err = decoder.Decode(param)
				// case "EraIndex":
				// 	param := types.EraIndex{}
				// 	err = decoder.Decode(param)
				// case "SessionIndex":
				// 	param := types.SessionIndex{}
				// 	err = decoder.Decode(param)
				// case "ElectionCompute":
				// 	param := types.ElectionCompute{}
				// 	err = decoder.Decode(param)
				// case "MemberCount":
				// 	param := types.MemberCount{}
				// 	err = decoder.Decode(param)
				// case "sp_std":
				// 	param := // types.sp_std::marker::PhantomData<(AccountId, Event)>{}
				// 	err = decoder.Decode(param)
				// case "Vec":
				// 	param := // types.Vec<(OracleKey, OracleValue)>{}
				// 	err = decoder.Decode(param)
				// case "CurrencyId":
				// 	param := types.CurrencyId{}
				// 	err = decoder.Decode(param)
				// case "Amount":
				// 	param := types.Amount{}
				// 	err = decoder.Decode(param)
				// case "VestingSchedule":
				// 	param := types.VestingSchedule{}
				// 	err = decoder.Decode(param)
				case "BlockNumber":
					param := types.BlockNumber(0)
					err = decoder.Decode(param)
				// case "DispatchId":
				// 	param := types.DispatchId{}
				// 	err = decoder.Decode(param)
				case "StorageKey":
					param := types.StorageKey{}
					err = decoder.Decode(param)
				// case "StorageValue":
				// 	param := types.StorageValue{}
				// 	err = decoder.Decode(param)
				// case "AuctionId":
				// 	param := types.AuctionId{}
				// 	err = decoder.Decode(param)
				// case "Price":
				// 	param := types.Price{}
				// 	err = decoder.Decode(param)
				// case "DebitAmount":
				// 	param := types.DebitAmount{}
				// 	err = decoder.Decode(param)
				// case "DebitBalance":
				// 	param := types.DebitBalance{}
				// 	err = decoder.Decode(param)
				// case "Share":
				// 	param := types.Share{}
				// 	err = decoder.Decode(param)
				// case "LiquidationStrategy":
				// 	param := types.LiquidationStrategy{}
				// 	err = decoder.Decode(param)
				// case "Option":
				// 	param := types.Option<Rate>{}
				// 	err = decoder.Decode(param)
				// case "Option":
				// 	param := types.Option<Ratio>{}
				// 	err = decoder.Decode(param)
				// case "Rate":
				// 	param := types.Rate{}
				// 	err = decoder.Decode(param)
				// case "Vec":
				// 	param := // types.Vec<(CurrencyId, Balance)>{}
				// 	err = decoder.Decode(param)
				// case "AirDropCurrencyId":
				// 	param := types.AirDropCurrencyId{}
				// 	err = decoder.Decode(param)
				// case "Vec":
				// 	param := // types.Vec<u8>{}
				// 	err = decoder.Decode(param)

				case "AuthorityList":
					param := []struct {
						AuthorityID     types.AuthorityID
						AuthorityWeight types.U64
					}{}
					err = decoder.Decode(param)
				default:
					return fmt.Errorf("unable to decode field %v_%v arg #%v %v", moduleName,
						eventName, j, eventParams.Args[j])
				}
			}

			fmt.Printf("unable to find field %v_%v for event #%v with EventID %v\n", moduleName, eventName, i, id)
			continue
			// return fmt.Errorf("unable to find field %v_%v for event #%v with EventID %v", moduleName, eventName, i, id)
		}

		// create a pointer to with the correct type that will hold the decoded event
		holder := reflect.New(field.Type().Elem())

		// ensure first field is for Phase, last field is for Topics
		numFields := holder.Elem().NumField()
		fmt.Printf("numFields: %v\n", numFields)
		if numFields < 2 {
			return fmt.Errorf("expected event #%v with EventID %v, field %v_%v to have at least 2 fields "+
				"(for Phase and Topics), but has %v fields", i, id, moduleName, eventName, numFields)
		}
		phaseField := holder.Elem().FieldByIndex([]int{0})
		if phaseField.Type() != reflect.TypeOf(phase) {
			return fmt.Errorf("expected the first field of event #%v with EventID %v, field %v_%v to be of type "+
				"types.Phase, but got %v", i, id, moduleName, eventName, phaseField.Type())
		}
		topicsField := holder.Elem().FieldByIndex([]int{numFields - 1})
		if topicsField.Type() != reflect.TypeOf([]types.Hash{}) {
			return fmt.Errorf("expected the last field of event #%v with EventID %v, field %v_%v to be of type "+
				"[]types.Hash for Topics, but got %v", i, id, moduleName, eventName, topicsField.Type())
		}

		// set the phase we decoded earlier
		phaseField.Set(reflect.ValueOf(phase))

		// set the remaining fields
		for j := 1; j < numFields; j++ {
			fmt.Printf("decoding field: %v\n", j)
			err = decoder.Decode(holder.Elem().FieldByIndex([]int{j}).Addr().Interface())
			if err != nil {
				return fmt.Errorf("unable to decode field %v event #%v with EventID %v, field %v_%v: %v", j, i, id, moduleName,
					eventName, err)
			}
		}

		// add the decoded event to the slice
		field.Set(reflect.Append(field, holder.Elem()))

		fmt.Println(fmt.Sprintf("decoded event #%v", i))
	}
	return nil
}

func findEventForEventID(m types.MetadataV10, eventID types.EventID) (*types.EventMetadataV4, error) {
	mi := uint8(0)
	for _, mod := range m.Modules {
		if !mod.HasEvents {
			continue
		}
		if mi != eventID[0] {
			mi++
			continue
		}
		if int(eventID[1]) >= len(mod.Events) {
			return nil, fmt.Errorf("event index %v for module %v out of range", eventID[1], mod.Name)
		}
		return &mod.Events[eventID[1]], nil
	}
	return nil, fmt.Errorf("module index %v out of range", eventID[0])
}
