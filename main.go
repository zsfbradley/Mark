package main

import (
	"log"
	"markV5/mError"
	"markV5/mFace"
	"markV5/mNet"
	"time"
)

func main() {
	server, err := mNet.DefaultServer()
	if pass := checkError(err); !pass {
		return
	}

	err = server.RegisterRoute("0000000001", route)
	if pass := checkError(err); !pass {
		return
	}
	err = server.RegisterRoutes(routes())
	if pass := checkError(err); !pass {
		return
	}

	err = server.RegisterFilter(mFace.MsgFilterType_Request, filter)
	if pass := checkError(err); !pass {
		return
	}
	err = server.RegisterFilter(mFace.MsgFilterType_Response, filter2)
	if pass := checkError(err); !pass {
		return
	}
	err = server.RegisterEntranceFunc(mainFunc)
	if pass := checkError(err); !pass {
		return
	}

	err = server.Start()
	if pass := checkError(err); !pass {
		return
	}

	time.Sleep(time.Second * time.Duration(10))

	err = server.Stop()
	if pass := checkError(err); !pass {
		return
	}

	for {
	}

	//time.Sleep(time.Duration(3)*time.Second)
	//
	//err = server.Reload()
	//if pass := checkError(err); !pass {
	//	return
	//}
	//
	//time.Sleep(time.Duration(3)*time.Second)
	//
	//err = server.Stop()
	//if pass := checkError(err); !pass {
	//	return
	//}
}

func route(request mFace.MMessage, response mFace.MMessage) mFace.MError {
	log.Printf("[route request] - %v", request)
	//log.Printf("[route response] - %v" , response)

	response.Unmarshal([]byte("曹尼玛"))

	return mError.Nil
}

func route2(request mFace.MMessage, response mFace.MMessage) mFace.MError {
	log.Printf("[route2 request] - %v", request)
	//log.Printf("[route response] - %v" , response)

	response.Unmarshal([]byte("曹尼lao玛"))
	return mError.Nil
}

func routes() map[string]mFace.RouteHandleFunc {
	routes := make(map[string]mFace.RouteHandleFunc)

	routes["0000000002"] = route
	routes["0000000003"] = route
	routes["0000000004"] = route
	routes["0000000005"] = route
	routes["0000000006"] = route2

	return routes
}

func filter(request mFace.MMessage) bool {
	//log.Printf("[filter request] - %v" , request)
	return true
}

func filter2(response mFace.MMessage) bool {
	//log.Printf("[filter response] - %v" , response)
	return true
}

func checkError(err mFace.MError) bool {
	if err.NotNil() {
		log.Println(err.Error())
		return false
	}
	return true
}

func mainFunc() mFace.MError {
	log.Println("Entrance")
	return mError.Nil
}
