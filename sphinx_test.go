package sphinx

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/brsuite/brond/btcec"
	"github.com/brsuite/brond/chaincfg"
	"github.com/davecgh/go-spew/spew"
)

// BOLT 4 Test Vectors
var (
	// bolt4PubKeys are the public keys of the hops used in the route.
	bolt4PubKeys = []string{
		"02eec7245d6b7d2ccb30380bfbe2a3648cd7a942653f5aa340edcea1f283686619",
		"0324653eac434488002cc06bbfb7f10fe18991e35f9fe4302dbea6d2353dc0ab1c",
		"027f31ebc5462c1fdce1b737ecff52d37d75dea43ce11c74d25aa297165faa2007",
		"032c0b7cf95324a07d05398b240174dc0c2be444d96b159aa6c7f7b1e668680991",
		"02edabbd16b41c8371b92ef2f04c1185b4f03b6dcd52ba9b78d9d7c89c8f221145",
	}

	// bolt4SessionKey is the session private key.
	bolt4SessionKey = bytes.Repeat([]byte{'A'}, 32)

	// bolt4AssocData is the associated data added to the packet.
	bolt4AssocData = bytes.Repeat([]byte{'B'}, 32)

	// bolt4FinalPacketHex encodes the expected sphinx packet as a result of
	// creating a new packet with the above parameters.
	bolt4FinalPacketHex = "0002eec7245d6b7d2ccb30380bfbe2a3648cd7" +
		"a942653f5aa340edcea1f283686619e5f14350c2a76fc232b5e4" +
		"6d421e9615471ab9e0bc887beff8c95fdb878f7b3a71da571226" +
		"458c510bbadd1276f045c21c520a07d35da256ef75b436796243" +
		"7b0dd10f7d61ab590531cf08000178a333a347f8b4072e216400" +
		"406bdf3bf038659793a86cae5f52d32f3438527b47a1cfc54285" +
		"a8afec3a4c9f3323db0c946f5d4cb2ce721caad69320c3a469a2" +
		"02f3e468c67eaf7a7cda226d0fd32f7b48084dca885d15222e60" +
		"826d5d971f64172d98e0760154400958f00e86697aa1aa9d41be" +
		"e8119a1ec866abe044a9ad635778ba61fc0776dc832b39451bd5" +
		"d35072d2269cf9b040d6ba38b54ec35f81d7fc67678c3be47274" +
		"f3c4cc472aff005c3469eb3bc140769ed4c7f0218ff8c6c7dd72" +
		"21d189c65b3b9aaa71a01484b122846c7c7b57e02e679ea8469b" +
		"70e14fe4f70fee4d87b910cf144be6fe48eef24da475c0b0bcc6" +
		"565ae82cd3f4e3b24c76eaa5616c6111343306ab35c1fe5ca4a7" +
		"7c0e314ed7dba39d6f1e0de791719c241a939cc493bea2bae1c1" +
		"e932679ea94d29084278513c77b899cc98059d06a27d171b0dbd" +
		"f6bee13ddc4fc17a0c4d2827d488436b57baa167544138ca2e64" +
		"a11b43ac8a06cd0c2fba2d4d900ed2d9205305e2d7383cc98dac" +
		"b078133de5f6fb6bed2ef26ba92cea28aafc3b9948dd9ae5559e" +
		"8bd6920b8cea462aa445ca6a95e0e7ba52961b181c79e73bd581" +
		"821df2b10173727a810c92b83b5ba4a0403eb710d2ca10689a35" +
		"bec6c3a708e9e92f7d78ff3c5d9989574b00c6736f84c199256e" +
		"76e19e78f0c98a9d580b4a658c84fc8f2096c2fbea8f5f8c59d0" +
		"fdacb3be2802ef802abbecb3aba4acaac69a0e965abd8981e989" +
		"6b1f6ef9d60f7a164b371af869fd0e48073742825e9434fc54da" +
		"837e120266d53302954843538ea7c6c3dbfb4ff3b2fdbe244437" +
		"f2a153ccf7bdb4c92aa08102d4f3cff2ae5ef86fab4653595e6a" +
		"5837fa2f3e29f27a9cde5966843fb847a4a61f1e76c281fe8bb2" +
		"b0a181d096100db5a1a5ce7a910238251a43ca556712eaadea16" +
		"7fb4d7d75825e440f3ecd782036d7574df8bceacb397abefc5f5" +
		"254d2722215c53ff54af8299aaaad642c6d72a14d27882d9bbd5" +
		"39e1cc7a527526ba89b8c037ad09120e98ab042d3e8652b31ae0" +
		"e478516bfaf88efca9f3676ffe99d2819dcaeb7610a626695f53" +
		"117665d267d3f7abebd6bbd6733f645c72c389f03855bdf1e4b8" +
		"075b516569b118233a0f0971d24b83113c0b096f5216a207ca99" +
		"a7cddc81c130923fe3d91e7508c9ac5f2e914ff5dccab9e55856" +
		"6fa14efb34ac98d878580814b94b73acbfde9072f30b881f7f0f" +
		"ff42d4045d1ace6322d86a97d164aa84d93a60498065cc7c20e6" +
		"36f5862dc81531a88c60305a2e59a985be327a6902e4bed986db" +
		"f4a0b50c217af0ea7fdf9ab37f9ea1a1aaa72f54cf40154ea9b2" +
		"69f1a7c09f9f43245109431a175d50e2db0132337baa0ef97eed" +
		"0fcf20489da36b79a1172faccc2f7ded7c60e00694282d93359c" +
		"4682135642bc81f433574aa8ef0c97b4ade7ca372c5ffc23c7ed" +
		"dd839bab4e0f14d6df15c9dbeab176bec8b5701cf054eb3072f6" +
		"dadc98f88819042bf10c407516ee58bce33fbe3b3d86a54255e5" +
		"77db4598e30a135361528c101683a5fcde7e8ba53f3456254be8" +
		"f45fe3a56120ae96ea3773631fcb3873aa3abd91bcff00bd38bd" +
		"43697a2e789e00da6077482e7b1b1a677b5afae4c54e6cbdf737" +
		"7b694eb7d7a5b913476a5be923322d3de06060fd5e819635232a" +
		"2cf4f0731da13b8546d1d6d4f8d75b9fce6c2341a71b0ea6f780" +
		"df54bfdb0dd5cd9855179f602f917265f21f9190c70217774a6f" +
		"baaa7d63ad64199f4664813b955cff954949076dcf"

	testLegacyRouteNumHops = 20
)

func newTestRoute(numHops int) ([]*Router, *PaymentPath, *[]HopData, *OnionPacket, error) {
	nodes := make([]*Router, numHops)

	// Create numHops random sphinx nodes.
	for i := 0; i < len(nodes); i++ {
		privKey, err := btcec.NewPrivateKey(btcec.S256())
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("Unable to "+
				"generate random key for sphinx node: %v", err)
		}

		nodes[i] = NewRouter(
			privKey, &chaincfg.MainNetParams, NewMemoryReplayLog(),
		)
	}

	// Gather all the pub keys in the path.
	var (
		route PaymentPath
	)
	for i := 0; i < len(nodes); i++ {
		hopData := HopData{
			ForwardAmount: uint64(i),
			OutgoingCltv:  uint32(i),
		}
		copy(hopData.NextAddress[:], bytes.Repeat([]byte{byte(i)}, 8))

		hopPayload, err := NewHopPayload(&hopData, nil)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("unable to "+
				"create new hop payload: %v", err)
		}

		route[i] = OnionHop{
			NodePub:    *nodes[i].onionKey.PubKey(),
			HopPayload: hopPayload,
		}
	}

	// Generate a forwarding message to route to the final node via the
	// generated intermediate nodes above.  Destination should be Hash160,
	// adding padding so parsing still works.
	sessionKey, _ := btcec.PrivKeyFromBytes(
		btcec.S256(), bytes.Repeat([]byte{'A'}, 32),
	)
	fwdMsg, err := NewOnionPacket(
		&route, sessionKey, nil, DeterministicPacketFiller,
	)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("unable to create "+
			"forwarding message: %#v", err)
	}

	var hopsData []HopData
	for i := 0; i < len(nodes); i++ {
		hopData, err := route[i].HopPayload.HopData()
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("unable to "+
				"gen hop data: %v", err)
		}

		hopsData = append(hopsData, *hopData)
	}

	return nodes, &route, &hopsData, fwdMsg, nil
}

func TestBolt4Packet(t *testing.T) {
	var (
		route    PaymentPath
		hopsData []HopData
	)
	for i, pubKeyHex := range bolt4PubKeys {
		pubKeyBytes, err := hex.DecodeString(pubKeyHex)
		if err != nil {
			t.Fatalf("unable to decode BOLT 4 hex pubkey #%d: %v", i, err)
		}

		pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
		if err != nil {
			t.Fatalf("unable to parse BOLT 4 pubkey #%d: %v", i, err)
		}

		hopData := HopData{
			ForwardAmount: uint64(i),
			OutgoingCltv:  uint32(i),
		}
		copy(hopData.NextAddress[:], bytes.Repeat([]byte{byte(i)}, 8))
		hopsData = append(hopsData, hopData)

		hopPayload, err := NewHopPayload(&hopData, nil)
		if err != nil {
			t.Fatalf("unable to make hop payload: %v", err)
		}

		pubKey.Curve = nil

		route[i] = OnionHop{
			NodePub:    *pubKey,
			HopPayload: hopPayload,
		}
	}

	finalPacket, err := hex.DecodeString(bolt4FinalPacketHex)
	if err != nil {
		t.Fatalf("unable to decode BOLT 4 final onion packet from hex: "+
			"%v", err)
	}

	sessionKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), bolt4SessionKey)
	pkt, err := NewOnionPacket(&route, sessionKey, bolt4AssocData, BlankPacketFiller)
	if err != nil {
		t.Fatalf("unable to construct onion packet: %v", err)
	}

	var b bytes.Buffer
	if err := pkt.Encode(&b); err != nil {
		t.Fatalf("unable to decode onion packet: %v", err)
	}

	if bytes.Compare(b.Bytes(), finalPacket) != 0 {
		t.Fatalf("final packet does not match expected BOLT 4 packet, "+
			"want: %s, got %s", hex.EncodeToString(finalPacket),
			hex.EncodeToString(b.Bytes()))
	}
}

func TestSphinxCorrectness(t *testing.T) {
	nodes, _, hopDatas, fwdMsg, err := newTestRoute(testLegacyRouteNumHops)
	if err != nil {
		t.Fatalf("unable to create random onion packet: %v", err)
	}

	// Now simulate the message propagating through the mix net eventually
	// reaching the final destination.
	for i := 0; i < len(nodes); i++ {
		// Start each node's ReplayLog and defer shutdown
		nodes[i].log.Start()
		defer nodes[i].log.Stop()

		hop := nodes[i]

		t.Logf("Processing at hop: %v \n", i)
		onionPacket, err := hop.ProcessOnionPacket(fwdMsg, nil, uint32(i)+1)
		if err != nil {
			t.Fatalf("Node %v was unable to process the "+
				"forwarding message: %v", i, err)
		}

		// The hop data for this hop should *exactly* match what was
		// initially used to construct the packet.
		expectedHopData := (*hopDatas)[i]
		if !reflect.DeepEqual(*onionPacket.ForwardingInstructions, expectedHopData) {
			t.Fatalf("hop data doesn't match: expected %v, got %v",
				spew.Sdump(expectedHopData),
				spew.Sdump(onionPacket.ForwardingInstructions))
		}

		// If this is the last hop on the path, the node should
		// recognize that it's the exit node.
		if i == len(nodes)-1 {
			if onionPacket.Action != ExitNode {
				t.Fatalf("Processing error, node %v is the last hop in "+
					"the path, yet it doesn't recognize so", i)
			}

		} else {
			// If this isn't the last node in the path, then the
			// returned action should indicate that there are more
			// hops to go.
			if onionPacket.Action != MoreHops {
				t.Fatalf("Processing error, node %v is not the final"+
					" hop, yet thinks it is.", i)
			}

			// The next hop should have been parsed as node[i+1].
			parsedNextHop := onionPacket.ForwardingInstructions.NextAddress[:]
			expected := bytes.Repeat([]byte{byte(i)}, AddressSize)
			if !bytes.Equal(parsedNextHop, expected) {
				t.Fatalf("Processing error, next hop parsed incorrectly."+
					" next hop should be %v, was instead parsed as %v",
					hex.EncodeToString(nodes[i+1].nodeID[:]),
					hex.EncodeToString(parsedNextHop))
			}

			fwdMsg = onionPacket.NextPacket
		}
	}
}

func TestSphinxSingleHop(t *testing.T) {
	// We'd like to test the proper behavior of the correctness of onion
	// packet processing for "single-hop" payments which bare a full onion
	// packet.
	nodes, _, _, fwdMsg, err := newTestRoute(1)
	if err != nil {
		t.Fatalf("unable to create test route: %v", err)
	}

	// Start the ReplayLog and defer shutdown
	nodes[0].log.Start()
	defer nodes[0].log.Stop()

	// Simulating a direct single-hop payment, send the sphinx packet to
	// the destination node, making it process the packet fully.
	processedPacket, err := nodes[0].ProcessOnionPacket(fwdMsg, nil, 1)
	if err != nil {
		t.Fatalf("unable to process sphinx packet: %v", err)
	}

	// The destination node should detect that the packet is destined for
	// itself.
	if processedPacket.Action != ExitNode {
		t.Fatalf("processed action is correct, is %v should be %v",
			processedPacket.Action, ExitNode)
	}
}

func TestSphinxNodeRelpay(t *testing.T) {
	// We'd like to ensure that the sphinx node itself rejects all replayed
	// packets which share the same shared secret.
	nodes, _, _, fwdMsg, err := newTestRoute(testLegacyRouteNumHops)
	if err != nil {
		t.Fatalf("unable to create test route: %v", err)
	}

	// Start the ReplayLog and defer shutdown
	nodes[0].log.Start()
	defer nodes[0].log.Stop()

	// Allow the node to process the initial packet, this should proceed
	// without any failures.
	if _, err := nodes[0].ProcessOnionPacket(fwdMsg, nil, 1); err != nil {
		t.Fatalf("unable to process sphinx packet: %v", err)
	}

	// Now, force the node to process the packet a second time, this should
	// fail with a detected replay error.
	if _, err := nodes[0].ProcessOnionPacket(fwdMsg, nil, 1); err != ErrReplayedPacket {
		t.Fatalf("sphinx packet replay should be rejected, instead error is %v", err)
	}
}

func TestSphinxNodeRelpaySameBatch(t *testing.T) {
	// We'd like to ensure that the sphinx node itself rejects all replayed
	// packets which share the same shared secret.
	nodes, _, _, fwdMsg, err := newTestRoute(testLegacyRouteNumHops)
	if err != nil {
		t.Fatalf("unable to create test route: %v", err)
	}

	// Start the ReplayLog and defer shutdown
	nodes[0].log.Start()
	defer nodes[0].log.Stop()

	tx := nodes[0].BeginTxn([]byte("0"), 2)

	// Allow the node to process the initial packet, this should proceed
	// without any failures.
	if err := tx.ProcessOnionPacket(0, fwdMsg, nil, 1); err != nil {
		t.Fatalf("unable to process sphinx packet: %v", err)
	}

	// Now, force the node to process the packet a second time, this call
	// should not fail, even though the batch has internally recorded this
	// as a duplicate.
	err = tx.ProcessOnionPacket(1, fwdMsg, nil, 1)
	if err != nil {
		t.Fatalf("adding duplicate sphinx packet to batch should not "+
			"result in an error, instead got: %v", err)
	}

	// Commit the batch to disk, then we will inspect the replay set to
	// ensure the duplicate entry was properly included.
	_, replaySet, err := tx.Commit()
	if err != nil {
		t.Fatalf("unable to commit batch of sphinx packets: %v", err)
	}

	if replaySet.Contains(0) {
		t.Fatalf("index 0 was not expected to be in replay set")
	}

	if !replaySet.Contains(1) {
		t.Fatalf("expected replay set to contain duplicate packet " +
			"at index 1")
	}
}

func TestSphinxNodeRelpayLaterBatch(t *testing.T) {
	// We'd like to ensure that the sphinx node itself rejects all replayed
	// packets which share the same shared secret.
	nodes, _, _, fwdMsg, err := newTestRoute(testLegacyRouteNumHops)
	if err != nil {
		t.Fatalf("unable to create test route: %v", err)
	}

	// Start the ReplayLog and defer shutdown
	nodes[0].log.Start()
	defer nodes[0].log.Stop()

	tx := nodes[0].BeginTxn([]byte("0"), 1)

	// Allow the node to process the initial packet, this should proceed
	// without any failures.
	if err := tx.ProcessOnionPacket(uint16(0), fwdMsg, nil, 1); err != nil {
		t.Fatalf("unable to process sphinx packet: %v", err)
	}

	_, _, err = tx.Commit()
	if err != nil {
		t.Fatalf("unable to commit sphinx batch: %v", err)
	}

	tx2 := nodes[0].BeginTxn([]byte("1"), 1)

	// Now, force the node to process the packet a second time, this should
	// fail with a detected replay error.
	err = tx2.ProcessOnionPacket(uint16(0), fwdMsg, nil, 1)
	if err != nil {
		t.Fatalf("sphinx packet replay should not have been rejected, "+
			"instead error is %v", err)
	}

	_, replays, err := tx2.Commit()
	if err != nil {
		t.Fatalf("unable to commit second sphinx batch: %v", err)
	}

	if !replays.Contains(0) {
		t.Fatalf("expected replay set to contain index: %v", 0)
	}
}

func TestSphinxNodeReplayBatchIdempotency(t *testing.T) {
	// We'd like to ensure that the sphinx node itself rejects all replayed
	// packets which share the same shared secret.
	nodes, _, _, fwdMsg, err := newTestRoute(testLegacyRouteNumHops)
	if err != nil {
		t.Fatalf("unable to create test route: %v", err)
	}

	// Start the ReplayLog and defer shutdown
	nodes[0].log.Start()
	defer nodes[0].log.Stop()

	tx := nodes[0].BeginTxn([]byte("0"), 1)

	// Allow the node to process the initial packet, this should proceed
	// without any failures.
	if err := tx.ProcessOnionPacket(uint16(0), fwdMsg, nil, 1); err != nil {
		t.Fatalf("unable to process sphinx packet: %v", err)
	}

	packets, replays, err := tx.Commit()
	if err != nil {
		t.Fatalf("unable to commit sphinx batch: %v", err)
	}

	tx2 := nodes[0].BeginTxn([]byte("0"), 1)

	// Now, force the node to process the packet a second time, this should
	// not fail with a detected replay error.
	err = tx2.ProcessOnionPacket(uint16(0), fwdMsg, nil, 1)
	if err != nil {
		t.Fatalf("sphinx packet replay should not have been rejected, "+
			"instead error is %v", err)
	}

	packets2, replays2, err := tx2.Commit()
	if err != nil {
		t.Fatalf("unable to commit second sphinx batch: %v", err)
	}

	if replays.Size() != replays2.Size() {
		t.Fatalf("expected replay set to be %v, instead got %v",
			replays, replays2)
	}

	if !reflect.DeepEqual(packets, packets2) {
		t.Fatalf("expected packets to be %v, instead go %v",
			packets, packets2)
	}
}

func TestSphinxAssocData(t *testing.T) {
	// We want to make sure that the associated data is considered in the
	// HMAC creation
	nodes, _, _, fwdMsg, err := newTestRoute(5)
	if err != nil {
		t.Fatalf("unable to create random onion packet: %v", err)
	}

	// Start the ReplayLog and defer shutdown
	nodes[0].log.Start()
	defer nodes[0].log.Stop()

	_, err = nodes[0].ProcessOnionPacket(fwdMsg, []byte("somethingelse"), 1)
	if err == nil {
		t.Fatalf("we should fail when associated data changes")
	}

}

func TestSphinxEncodeDecode(t *testing.T) {
	// Create some test data with a randomly populated, yet valid onion
	// forwarding message.
	_, _, _, fwdMsg, err := newTestRoute(5)
	if err != nil {
		t.Fatalf("unable to create random onion packet: %v", err)
	}

	// Encode the created onion packet into an empty buffer. This should
	// succeeed without any errors.
	var b bytes.Buffer
	if err := fwdMsg.Encode(&b); err != nil {
		t.Fatalf("unable to encode message: %v", err)
	}

	// Now decode the bytes encoded above. Again, this should succeeed
	// without any errors.
	newFwdMsg := &OnionPacket{}
	if err := newFwdMsg.Decode(&b); err != nil {
		t.Fatalf("unable to decode message: %v", err)
	}

	// The two forwarding messages should now be identical.
	if !reflect.DeepEqual(fwdMsg, newFwdMsg) {
		t.Fatalf("forwarding messages don't match, %v vs %v",
			spew.Sdump(fwdMsg), spew.Sdump(newFwdMsg))
	}
}

func newEOBRoute(numHops uint32,
	eobMapping map[int]HopPayload) (*OnionPacket, []*Router, error) {

	nodes := make([]*Router, numHops)

	if uint32(len(eobMapping)) != numHops {
		return nil, nil, fmt.Errorf("must provide payload " +
			"mapping for all hops")
	}

	// First, we'll assemble a set of routers that will consume all the
	// hops we create in this path.
	for i := 0; i < len(nodes); i++ {
		privKey, err := btcec.NewPrivateKey(btcec.S256())
		if err != nil {
			return nil, nil, fmt.Errorf("Unable to generate "+
				"random key for sphinx node: %v", err)
		}

		nodes[i] = NewRouter(
			privKey, &chaincfg.MainNetParams, NewMemoryReplayLog(),
		)
	}

	// Next we'll gather all the pubkeys in the path, checking our eob
	// mapping to see which hops need an extra payload.
	var (
		route PaymentPath
	)
	for i := 0; i < len(nodes); i++ {
		route[i] = OnionHop{
			NodePub:    *nodes[i].onionKey.PubKey(),
			HopPayload: eobMapping[i],
		}
	}

	// Generate a forwarding message to route to the final node via the
	// generated intermediate nodes above.  Destination should be Hash160,
	// adding padding so parsing still works.
	sessionKey, _ := btcec.PrivKeyFromBytes(
		btcec.S256(), bytes.Repeat([]byte{'A'}, 32),
	)
	fwdMsg, err := NewOnionPacket(
		&route, sessionKey, nil, DeterministicPacketFiller,
	)
	if err != nil {
		return nil, nil, err
	}

	return fwdMsg, nodes, nil
}

func mustNewHopPayload(hopData *HopData, eob []byte) HopPayload {
	payload, err := NewHopPayload(hopData, eob)
	if err != nil {
		panic(err)
	}

	return payload
}

// TestSphinxHopVariableSizedPayloads tests that we're able to fully decode an
// EOB payload that was targeted at the final hop in a route, and also when
// intermediate nodes have EOB data encoded as well. Additionally, we test that
// we're able to mix the legacy and current format within the same route.
func TestSphinxHopVariableSizedPayloads(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		numNodes      uint32
		eobMapping    map[int]HopPayload
		expectedError error
	}{
		// A single hop route with a payload going to the last hop in
		// the route. The payload is enough to fit into what would be
		// the normal frame type, but it's a TLV hop.
		{
			numNodes: 1,
			eobMapping: map[int]HopPayload{
				0: HopPayload{
					Type:    PayloadTLV,
					Payload: bytes.Repeat([]byte("a"), LegacyHopDataSize-HMACSize),
				},
			},
		},

		// A single hop route where the payload to the final node needs
		// to shift more than a single frame.
		{
			numNodes: 1,
			eobMapping: map[int]HopPayload{
				0: HopPayload{
					Type:    PayloadTLV,
					Payload: bytes.Repeat([]byte("a"), LegacyHopDataSize*3),
				},
			},
		},

		// A two hop route, so one going over 3 nodes, with the sender
		// encrypting a payload to the final node. The payload of the
		// final node will require more shifts than normal to parse the
		// data The first hop is a legacy hop containing the usual
		// amount of data.
		{
			numNodes: 2,
			eobMapping: map[int]HopPayload{
				0: mustNewHopPayload(&HopData{
					Realm:         [1]byte{0x00},
					ForwardAmount: 2,
					OutgoingCltv:  3,
					NextAddress:   [8]byte{1, 1, 1, 1, 1, 1, 1, 1},
				}, nil),
				1: HopPayload{
					Type:    PayloadTLV,
					Payload: bytes.Repeat([]byte("a"), LegacyHopDataSize*2),
				},
			},
		},

		// A 3 hop route (4 nodes) with all but the middle node
		// receiving a TLV payload. Each of the TLV hops will use a
		// distinct amount of data in each hop.
		{
			numNodes: 3,
			eobMapping: map[int]HopPayload{
				0: HopPayload{
					Type:    PayloadTLV,
					Payload: bytes.Repeat([]byte("a"), 100),
				},
				1: mustNewHopPayload(&HopData{
					Realm:         [1]byte{0x00},
					ForwardAmount: 22,
					OutgoingCltv:  9,
					NextAddress:   [8]byte{1, 1, 1, 1, 1, 1, 1, 1},
				}, nil),
				2: HopPayload{
					Type:    PayloadTLV,
					Payload: bytes.Repeat([]byte("a"), 256),
				},
			},
		},

		// A 3 hop route (4 nodes), each hop is a TLV hop and will use
		// a distinct amount of data for each of their hops.
		{
			numNodes: 3,
			eobMapping: map[int]HopPayload{
				0: HopPayload{
					Type:    PayloadTLV,
					Payload: bytes.Repeat([]byte("a"), 200),
				},
				1: HopPayload{
					Type:    PayloadTLV,
					Payload: bytes.Repeat([]byte("a"), 256),
				},
				2: HopPayload{
					Type:    PayloadTLV,
					Payload: bytes.Repeat([]byte("a"), 150),
				},
			},
		},

		// A 3 hop route (4 nodes) that carries more data then what fits
		// in the routing info.
		{
			numNodes: 3,
			eobMapping: map[int]HopPayload{
				0: HopPayload{
					Type:    PayloadTLV,
					Payload: bytes.Repeat([]byte("a"), 500),
				},
				1: HopPayload{
					Type:    PayloadTLV,
					Payload: bytes.Repeat([]byte("a"), 500),
				},
				2: HopPayload{
					Type:    PayloadTLV,
					Payload: bytes.Repeat([]byte("a"), 500),
				},
			},
			expectedError: ErrMaxRoutingInfoSizeExceeded,
		},
	}

	for testCaseNum, testCase := range testCases {
		nextPkt, routers, err := newEOBRoute(
			testCase.numNodes, testCase.eobMapping,
		)
		if testCase.expectedError != err {
			t.Fatalf("#%v: unable to create eob "+
				"route: %v", testCase, err)
		}
		if err != nil {
			continue
		}

		// We'll now walk thru manually each actual hop within the
		// route. We use the size of the routers rather than the number
		// of hops here as virtual EOB hops may have been inserted into
		// the route.
		for i := 0; i < len(routers); i++ {
			// Start each node's ReplayLog and defer shutdown
			routers[i].log.Start()
			defer routers[i].log.Stop()

			currentHop := routers[i]

			// Ensure that this hop is able to properly process
			// this onion packet. If additional EOB hops were
			// added, then it should be able to properly decrypt
			// all the layers and pass them on to the next node
			// properly.
			processedPacket, err := currentHop.ProcessOnionPacket(
				nextPkt, nil, uint32(i),
			)
			if err != nil {
				t.Fatalf("#%v: unable to process packet at "+
					"hop #%v: %v", testCaseNum, i, err)
			}

			// If this hop is expected to have EOB data, then we'll
			// check now to ensure the bytes were properly
			// recovered on the other end.
			eobData := testCase.eobMapping[i]
			if !reflect.DeepEqual(eobData.Payload,
				processedPacket.Payload.Payload) {
				t.Fatalf("#%v (hop %v): eob mismatch: expected "+
					"%v, got %v", testCaseNum, i,
					spew.Sdump(eobData.Payload),
					spew.Sdump(processedPacket.Payload.Payload))
			}

			if eobData.Type != processedPacket.Payload.Type {
				t.Fatalf("mismatched types: expected %v "+
					"got %v", eobData.Type,
					processedPacket.Payload.Type)
			}

			// If this is the last node (but not necessarily hop
			// due to EOB expansion), then it should recognize that
			// it's the exit node.
			if i == len(routers)-1 {
				if processedPacket.Action != ExitNode {
					t.Fatalf("#%v: Processing error, "+
						"node %v is the last hop in "+
						"the path, yet it doesn't "+
						"recognize so", testCaseNum, i)
				}
				continue
			}

			// If this isn't the last node in the path, then the
			// returned action should indicate that there are more
			// hops to go.
			if processedPacket.Action != MoreHops {
				t.Fatalf("#%v: Processing error, node %v is "+
					"not the final hop, yet thinks it is.",
					testCaseNum, i)
			}

			// The next hop should have been parsed as node[i+1],
			// but only if this was a legacy hop.
			if processedPacket.ForwardingInstructions != nil {
				parsedNextHop := processedPacket.ForwardingInstructions.NextAddress[:]

				expected := bytes.Repeat([]byte{byte(1)}, AddressSize)
				if !bytes.Equal(parsedNextHop, expected) {
					t.Fatalf("#%v: Processing error, next hop parsed "+
						"incorrectly. next hop should be %v, "+
						"was instead parsed as %v", testCaseNum,
						hex.EncodeToString(expected),
						hex.EncodeToString(parsedNextHop))
				}
			}

			nextPkt = processedPacket.NextPacket
		}
	}
}

// testFileName is the name of the multi-frame onion test file.
const testFileName = "testdata/onion-test-multi-frame.json"

type jsonHop struct {
	Type string `json:"type"`

	Pubkey string `json:"pubkey"`

	Payload string `json:"payload"`
}

type payloadTestCase struct {
	SessionKey string `json:"session_key"`

	AssociatedData string `json:"associated_data"`

	Hops []jsonHop `json:"hops"`
}

type jsonTestCase struct {
	Comment string `json:"comment"`

	Generate payloadTestCase `json:"generate"`

	Onion string `json:"onion"`

	Decode []string `json:"decode"`
}

// jsonTypeToPayloadType maps the JSON payload type to our concrete PayloadType
// type.
func jsonTypeToPayloadType(jsonType string) PayloadType {
	switch jsonType {
	case "raw":
		fallthrough
	case "tlv":
		return PayloadTLV

	case "legacy":
		return PayloadLegacy

	default:
		panic(fmt.Sprintf("unknown payload type: %v", jsonType))
	}
}

// TestVariablePayloadOnion tests that if we construct a packet that contains a
// mix of the old and new payload format, that we match the version that's
// included in the spec.
func TestVariablePayloadOnion(t *testing.T) {
	t.Parallel()

	// First, we'll read out the raw JSOn file at the target location.
	jsonBytes, err := ioutil.ReadFile(testFileName)
	if err != nil {
		t.Fatalf("unable to read json file: %v", err)
	}

	// Once we have the raw file, we'll unpack it into our jsonTestCase
	// struct defined above.
	testCase := &jsonTestCase{}
	if err := json.Unmarshal(jsonBytes, testCase); err != nil {
		t.Fatalf("unable to parse spec json file: %v", err)
	}

	// Next, we'll populate a new OnionHop using the information included
	// in this test case.
	var route PaymentPath
	for i, hop := range testCase.Generate.Hops {
		pubKeyBytes, err := hex.DecodeString(hop.Pubkey)
		if err != nil {
			t.Fatalf("unable to decode pubkey: %v", err)
		}
		pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
		if err != nil {
			t.Fatalf("unable to parse BOLT 4 pubkey #%d: %v", i, err)
		}

		payload, err := hex.DecodeString(hop.Payload)
		if err != nil {
			t.Fatalf("unable to decode payload: %v", err)
		}

		payloadType := jsonTypeToPayloadType(hop.Type)
		route[i] = OnionHop{
			NodePub: *pubKey,
			HopPayload: HopPayload{
				Type:    payloadType,
				Payload: payload,
			},
		}

		if payloadType == PayloadLegacy {
			route[i].HopPayload.Payload = append(
				[]byte{0x00}, route[i].HopPayload.Payload...,
			)

			route[i].HopPayload.Payload = append(
				route[i].HopPayload.Payload,
				bytes.Repeat([]byte{0x00}, NumPaddingBytes)...,
			)
		}
	}

	finalPacket, err := hex.DecodeString(testCase.Onion)
	if err != nil {
		t.Fatalf("unable to decode packet: %v", err)
	}

	sessionKeyBytes, err := hex.DecodeString(testCase.Generate.SessionKey)
	if err != nil {
		t.Fatalf("unable to generate session key: %v", err)
	}

	associatedData, err := hex.DecodeString(testCase.Generate.AssociatedData)
	if err != nil {
		t.Fatalf("unable to decode AD: %v", err)
	}

	// With all the required data assembled, we'll craft a new packet.
	sessionKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), sessionKeyBytes)
	pkt, err := NewOnionPacket(
		&route, sessionKey, associatedData, BlankPacketFiller,
	)
	if err != nil {
		t.Fatalf("unable to construct onion packet: %v", err)
	}

	var b bytes.Buffer
	if err := pkt.Encode(&b); err != nil {
		t.Fatalf("unable to decode onion packet: %v", err)
	}

	// Finally, we expect that our packet matches the packet included in
	// the spec's test vectors.
	if bytes.Compare(b.Bytes(), finalPacket) != 0 {
		t.Fatalf("final packet does not match expected BOLT 4 packet, "+
			"want: %s, got %s", hex.EncodeToString(finalPacket),
			hex.EncodeToString(b.Bytes()))
	}
}
