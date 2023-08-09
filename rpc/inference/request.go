package inference

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/**
Proto Package Installation:
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

Protoc generation:
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    inference.proto
*/

type InferenceTx struct {
	hash   string
	model  string
	params string
}

type InferenceConsolidation struct {
	tx           InferenceTx
	result       string
	attestations []string
	weight       float64
}

func (ic InferenceConsolidation) attest(threshold float64, node string, nodeWeight float64) bool {
	ic.attestations = append(ic.attestations, node)
	ic.weight += nodeWeight
	return ic.weight >= threshold
}

type RequestClient struct {
	port  int
	txMap map[string]float64
}

// Instantiating a new request client
func NewRequestClient(portNum int) *RequestClient {
	rc := &RequestClient{
		port: portNum,
	}
	return rc
}

// Emit inference transaction
func (rc RequestClient) emit(tx InferenceTx) (float64, error) {
	resultMap := make(map[string]InferenceConsolidation)
	nodes := getNodes()
	for _, ip := range nodes {
		serverAddr := getAddress(ip, rc.port)
		opts := getDialOptions()

		fmt.Println("Sending request to " + serverAddr)
		conn, err := grpc.Dial(serverAddr, opts...)
		if err != nil {
			log.Println("Unable to dial " + ip)
			continue
		}
		defer conn.Close()
		client := NewInferenceClient(conn)
		result := RunInference(client, tx)
		if _, ok := resultMap[result]; !ok {
			resultMap[result] = InferenceConsolidation{tx: tx, result: result, attestations: []string{}, weight: 0}
		}
		// Increment results count
		if val, ok := resultMap[result]; ok {
			complete := val.attest(float64(len(nodes)), ip, 1)
			if complete {
				parsedResult, err := strconv.ParseFloat(val.result, 64)
				if err == nil {
					return parsedResult, nil
				} else {
					return 0, err
				}
			}
		}
		/**
		TODO: ASYNCHRONOUS APPROACH FOR AFTER WE IMPLEMENT OPTIMISTIC PROCESSING
		go func(address string, port int) {
			// Endpoint configurations
			serverAddr := getAddress(ip, port)
			opts := getDialOptions()

			fmt.Println("Sending request to " + serverAddr)
			conn, err := grpc.Dial(serverAddr, opts...)
			if err != nil {
				fmt.Println(err)
			}
			defer conn.Close()
			client := NewInferenceClient(conn)
			result := RunInference(client, tx)
			// Increment results count
			if val, ok := resultMap[result]; ok {
				complete := val.attest(float64(len(nodes)), ip, 1)
				if complete {
					parsedResult, err := strconv.ParseFloat(val.result, 64)
					if (err != nil) {
						rc.txMap[tx.hash] = parsedResult
					} else {
						fmt.Println("Error encountered ")
					}
				}
			} else {
				resultMap[result] = InferenceConsolidation{tx: tx, result: result, attestations: []string{}}
			}

		}(ip, rc.port)
		*/
		// Hackathon iterative approach
	}
	return 0, errors.New("could not reach consensus")
}

// Runs inference request via gRPC
func RunInference(client InferenceClient, tx InferenceTx) string {
	inferenceParams := buildInferenceParameters(tx)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := client.RunInference(ctx, inferenceParams)
	if err != nil {
		log.Fatalf("RPC Failed: %v", err)
	}
	return result.Value
}

// Get IP addresses of inference nodes on network
func getNodes() []string {
	// TODO: DNS lookup of public IPV4 addresses of inference nodes
	// hard code for now
	nodes := []string{"3.139.238.241"}
	return nodes
}

func getWeightThreshold() float64 {
	return float64(len(getNodes()))
}

func getAddress(ip string, port int) string {
	return ip + ":" + strconv.Itoa(port)
}

func getDialOptions() []grpc.DialOption {
	var opts []grpc.DialOption
	// TODO: Add TLS and security auth measures
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return opts
}

func buildInferenceParameters(tx InferenceTx) *InferenceParameters {
	return &InferenceParameters{ModelHash: tx.model, ModelInput: tx.params}
}
