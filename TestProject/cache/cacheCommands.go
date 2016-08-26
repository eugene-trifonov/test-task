package cache

//Represents a Telnet-kind interface for a Cache conection.
//Shows allowed functionality.
type CacheCommands interface {

	//Returns a value that can be found in a Cache by passed <params>.
	//Or error if something goes wrong, e.g. <params> are wrong.
	GetValue(params []string) (interface{}, error)

	//Sets a key-value pair into a Cache. Where a value is a string.
	//Returns a value that was replaced, or <nil> there was no value for passed <params>.
	//Or error if something goes wrong, e.g. <params> are wrong.
	SetValue(params []string) (interface{}, error)

	//Updates an existing value in a Cache by passed <params>. A value is a string only.
	//Returns <true> if updated was successful, otherwise return <false>.
	//Or error if something goes wrong, e.g. <params> are wrong.
	UpdateValue(params []string) (bool, error)

	//Removes a value from a Cache by passed <params>.
	//Returns a removed value if any, otherwise <nil>.
	//Or error if something goes wrong, e.g. <params> are wrong.
	RemoveValue(params []string) (interface{}, bool, error)

	//Returns a value from a stored list in a Cache.
	//Or error if something goes wrong, e.g. <params> are wrong.
	GetListValue(params []string) (interface{}, error)

	//Appends a value to a list stored in a Cache.
	//Returns an error if something goes wrong, e.g. <params> are wrong.
	AppendListValue(params []string) error

	//Deletes a value from a stored list in a Cache.
	//Returns a removed value.
	//Or error if something goes wrong, e.g. <params> are wrong.
	DeleteListValue(params []string) (interface{}, error)

	//Returns a size of stored list in a Cache.
	//Or error if something goes wrong, e.g. <params> are wrong.
	GetListSize(params []string) (int, error)

	//Returns a value from a stored dictionary in a Cache.
	//Or error if something goes wrong, e.g. <params> are wrong.
	GetDictValue(params []string) (interface{}, error)

	//Sets a key-value pair into a dictionary stored in a Cache.
	//Returns a value that was repalced if any, otherwise <nil>.
	//Or error if something goes wrong, e.g. <params> are wrong.
	SetDictValue(params []string) (interface{}, error)

	//Appends a new key-value pair into a dictionary stored in a Cache if doesn't exist.
	//Returns <true> if a key-value pair was set up, otherwise returns <false>.
	//Or error if something goes wrong, e.g. <params> are wrong.
	AppendDictValue(params []string) (bool, error)

	//Deletes a key-value pair from a dictionary stored in a Cache.
	//Returns a removed value if exists, otherwise <nil>.
	//Or error if something goes wrong, e.g. <params> are wrong.
	DeleteDictValue(params []string) (interface{}, error)

	//Returns size of a dictionary stored in a Cache, a number of its key-value pairs.
	//Or error if something goes wrong, e.g. <params> are wrong.
	GetDictSize(params []string) (int, error)

	//Returns keys stored in a Cache.
	//Or error if something goes wrong, e.g. <params> are wrong.
	GetKeys(params []string) ([]string, error)

	//Updates a time to live of any value stored in a Cache.
	//Returns <true> if update was successful, otherwise returns <false>.
	//Or error if something goes wrong, e.g. <params> are wrong.
	UpdateTTL(params []string) (bool, error)

	//Returns size of a Cache, the count of its key-value pairs.
	//Or error if something goes wrong, e.g. <params> are wrong.
	GetSize() int
}
