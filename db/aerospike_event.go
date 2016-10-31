package db

import (
	"fmt"
	aerospike "github.com/aerospike/aerospike-client-go"
)

func NotDuplicate(hash string, ttl uint32, callback func()) {

	key, err := aerospike.NewKey("bidder", "duplicates", hash)

	if err != nil {
		fmt.Println(err)
		return
	}

	if result, err := client.Get(nil, key); err == nil {

		// fmt.Println("NotDuplicate: result", result)

		if result != nil {

			// fmt.Println("NOP")

		} else {

			policy := aerospike.NewWritePolicy(1, ttl)
			op := aerospike.AddOp(aerospike.NewBin("hit", 1))

			// record, err :=
			client.Operate(policy, key, op)
			// fmt.Println("ADD", record, err)

			callback()
		}
	}

	// var key = app.globals.aerospike.aerospike.key('bidder','duplicates',hash);
	// return app.globals.aerospike.db.get(key, function(getErr, rec, meta) {
	//     if(getErr.code==0)
	//     {
	//         return;
	//     }
	//     else
	//     {
	//         callback();
	//         var metadata = {ttl: timeout, gen: 1};
	//         var op = app.globals.aerospike.aerospike.operator;
	//         var ops = [ op.incr('hit', 1)];
	//         return app.globals.aerospike.db.operate(key,ops,metadata, function(err, key){
	//             return;
	//         });
	//     }
	// });
}
