package digibyte

import (
	"math/big"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

func init() {
	if err := chaincfg.Register(&MainNetParams); err != nil {
		panic(err)
	}
	if err := chaincfg.Register(&RegressionNetParams); err != nil {
		panic(err)
	}
}

var (
	bigOne       = big.NewInt(1)
	mainPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne)
)

const (
	DeploymentTestDummy = iota
	DeploymentCSV
	DeploymentSegwit
	// DeploymentVersionBits
	// DeploymentVersionReserveAlgos
	// DeploymentOdo
	// DeploymentEquihash
	// DeploymentEthash
	DefinedDeployments
)

// genesisCoinbaseTx is the coinbase transaction for the genesis blocks for
// the main network, regression test network, and test network (version 3).
var genesisCoinbaseTx = wire.MsgTx{
	Version: 1,
	TxIn: []*wire.TxIn{
		{
			PreviousOutPoint: wire.OutPoint{
				Hash:  chainhash.Hash{},
				Index: 0xffffffff,
			},
			SignatureScript: []byte{
				0x04, 0xff, 0xff, 0x00, 0x1d, 0x01, 0x04, 0x45, 0x55, 0x53, 0x41, 0x20, 0x54, 0x6f, 0x64, 0x61, /* |.......EUSA Toda| */
				0x79, 0x3a, 0x20, 0x31, 0x30, 0x2f, 0x4a, 0x61, 0x6e, 0x2f, 0x32, 0x30, 0x31, 0x34, 0x2c, 0x20, /* |y: 10/Jan/2014, | */
				0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x3a, 0x20, 0x44, 0x61, 0x74, 0x61, 0x20, 0x73, 0x74, 0x6f, /* |Target: Data sto| */
				0x6c, 0x65, 0x6e, 0x20, 0x66, 0x72, 0x6f, 0x6d, 0x20, 0x75, 0x70, 0x20, 0x74, 0x6f, 0x20, 0x31, /* |len from up to 1| */
				0x31, 0x30, 0x4d, 0x20, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x73, 0x61, 0x32, 0x30, /* |10M customers|    */
			},
			Sequence: 0xffffffff,
		},
	},
	TxOut: []*wire.TxOut{
		{
			Value: 0x12a05f200,
			PkScript: []byte{ // ToDo
				0x41, 0x04, 0x67, 0x8a, 0xfd, 0xb0, 0xfe, 0x55, /* |A.g....U| */
				0x48, 0x27, 0x19, 0x67, 0xf1, 0xa6, 0x71, 0x30, /* |H'.g..q0| */
				0xb7, 0x10, 0x5c, 0xd6, 0xa8, 0x28, 0xe0, 0x39, /* |..\..(.9| */
				0x09, 0xa6, 0x79, 0x62, 0xe0, 0xea, 0x1f, 0x61, /* |..yb...a| */
				0xde, 0xb6, 0x49, 0xf6, 0xbc, 0x3f, 0x4c, 0xef, /* |..I..?L.| */
				0x38, 0xc4, 0xf3, 0x55, 0x04, 0xe5, 0x1e, 0xc1, /* |8..U....| */
				0x12, 0xde, 0x5c, 0x38, 0x4d, 0xf7, 0xba, 0x0b, /* |..\8M...| */
				0x8d, 0x57, 0x8a, 0x4c, 0x70, 0x2b, 0x6b, 0xf1, /* |.W.Lp+k.| */
				0x1d, 0x5f, 0xac, /* |._.| */
			},
		},
	},
	LockTime: 0,
}

// USA Today: 10/Jan/2014, Target: Data stolen from up to 110M customers
var genesisMerkleRoot = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x96, 0x84, 0x1e, 0x6e, 0xcc, 0x8d,
	0xc9, 0x64, 0x3a, 0xad, 0xdf, 0xb6,
	0xfc, 0xd6, 0x16, 0xe0, 0x8f, 0x07,
	0x77, 0xc8, 0x7b, 0x50, 0x8f, 0x1c,
	0x9f, 0xb3, 0x5e, 0x46, 0x1b, 0xea,
	0x97, 0x74,
})

var genesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: genesisMerkleRoot,        // 7497ea1b465eb39f1c8f507bc877078fe016d6fcb6dfad3a64c98dcc6e1e8496
		Timestamp:  time.Unix(1389388394, 0), // 2014-01-10T21:13:14.000Z
		Bits:       0x1e0ffff0,               // 486604799 [00000000ffff0000000000000000000000000000000000000000000000000000]
		Nonce:      2447652,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

var genesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x96, 0x84, 0x1e, 0x6e, 0xcc, 0x8d, 0xc9, 0x64,
	0x3a, 0xad, 0xdf, 0xb6, 0xfc, 0xd6, 0x16, 0xe0,
	0x8f, 0x07, 0x77, 0xc8, 0x7b, 0x50, 0x8f, 0x1c,
	0x9f, 0xb3, 0x5e, 0x46, 0x1b, 0xea, 0x97, 0x74,
})

func newHashFromStr(hexStr string) *chainhash.Hash {
	hash, err := chainhash.NewHashFromStr(hexStr)
	if err != nil {
		panic(err)
	}
	return hash
}

var MainNetParams = chaincfg.Params{
	Name:        "mainnet",
	Net:         0xdab6c3fa,
	DefaultPort: "12024",
	DNSSeeds: []chaincfg.DNSSeed{
		{Host: "seed1.digibyte.io", HasFiltering: false},
		{Host: "seed2.digibyte.io", HasFiltering: false},
		{Host: "seed3.digibyte.io", HasFiltering: false},
		{Host: "seed.digibyte.io", HasFiltering: false},
		{Host: "digihash.co", HasFiltering: false},
		{Host: "digiexplorer.info", HasFiltering: false},
		{Host: "seed.digibyteprojects.com", HasFiltering: false},
	},

	// Chain parameters
	GenesisBlock:             &genesisBlock,
	GenesisHash:              &genesisHash,
	PowLimit:                 new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne),
	PowLimitBits:             0x1d00ffff,
	BIP0034Height:            4394880, // add8ca420f557f62377ec2be6e6f47b96cf2e68160d58aeb7b73433de834cca0
	BIP0065Height:            4394880, // add8ca420f557f62377ec2be6e6f47b96cf2e68160d58aeb7b73433de834cca0
	BIP0066Height:            4394880, // add8ca420f557f62377ec2be6e6f47b96cf2e68160d58aeb7b73433de834cca0
	CoinbaseMaturity:         100,
	SubsidyReductionInterval: 210000,
	TargetTimespan:           12 * time.Hour / 5, // 2.4 hours
	TargetTimePerBlock:       time.Second * 60,   // 60 seconds
	RetargetAdjustmentFactor: 4,                  // 25% less, 400% more
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0,
	GenerateSupported:        false,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []chaincfg.Checkpoint{
		{Height: 100000, Hash: newHashFromStr("236eb6c24599c6a6a2a3f2263086e1e67b2fe29a37bec19be89bf24b33f12cc9")},
		{Height: 200000, Hash: newHashFromStr("0000000000000e8acbba6f1d21394c5df111dd836a17b750f29150f72be02802")},
		{Height: 300000, Hash: newHashFromStr("701771ab8095e412122d1b7d49c60b0ef871751bc26a7d2a9a68d59dcff5edef")},
		{Height: 400000, Hash: newHashFromStr("596eeb60dab06afccb354f63abbd684e11fd46339977c7943528e4aa34ed0400")},
		{Height: 500000, Hash: newHashFromStr("745137430c905f1a82a59cd733832c543724043fb2fd60e532bac21b21256bb3")},
		{Height: 600000, Hash: newHashFromStr("292bee9bd10c2d02ebebc2a366f8828f79a8f2369c842fb090f8537e87b73487")},
		{Height: 700000, Hash: newHashFromStr("c68e02db2c1ef531d37158b2528bf57f6e1c8596d383e221a8ba4529bd02c087")},
		{Height: 800000, Hash: newHashFromStr("cbd4702a5eb54723edbae45659871449b300b3da9a30e5e4c6824fac80cae3e7")},
		{Height: 900000, Hash: newHashFromStr("11a7336b91bf2994321b8bcd49ab2fcbe45afa6399b406b59620915db0e11a11")},
		{Height: 1000000, Hash: newHashFromStr("b85849392739f50e1fc7568b16eeb1bf37f36c144d7ec46340a4c8a202253ade")},
		{Height: 1100000, Hash: newHashFromStr("b4b6536e012cf7f80a6e46d8801a5deec3106c8988272be8907f726e94c38407")},
		{Height: 1200000, Hash: newHashFromStr("000000000000083784693fa8994f64417fb2feb608bb2d0c93b60afeffdd48e3")},
		{Height: 1300000, Hash: newHashFromStr("0000000000000020ff988770a99cf91c9ddabf97122329a9cedc8c5d4d774a12")},
		{Height: 1400000, Hash: newHashFromStr("1cdc20a5e53a84c427d903eb7f748be09a0d3a9d8faf46ed6603857c0f424840")},
		{Height: 1500000, Hash: newHashFromStr("af353af43bb01a5168f7d344681de8389c44d2e431841a1943a1467f8d851a8b")},
		{Height: 1600000, Hash: newHashFromStr("98efd15a6deec563359366dad7dc180f4b8041aa42f9a7e5fd907152648131e3")},
		{Height: 1700000, Hash: newHashFromStr("000000000000069834ceddccc356b743970fdd07a7596a8d7bd3b30fcb11ba32")},
		{Height: 1800000, Hash: newHashFromStr("72f46e1fff56518dce7e540b407260ea827cb1c4652f24eb1d1917f54b95d65a")},
		{Height: 1900000, Hash: newHashFromStr("66d74c12a3b939fdec4f6f6bbe30a74a93bf5171989c264e19da00785b05eb8b")},
		{Height: 2000000, Hash: newHashFromStr("10f522ec60d8af2e2cbd9e2268260c33fb8bbf9cd9f176b4fddcae7493c6791d")},
		{Height: 2100000, Hash: newHashFromStr("f73ad3dd4ad9fe97d244758b9d564fe6a631c5bfd5e99b54b854c8a60ca93b03")},
		{Height: 2200000, Hash: newHashFromStr("f2cda55b21320d1d896cfbcd27f803bbcc682f671cc28ed45cf1e835f4f8c722")},
		{Height: 2300000, Hash: newHashFromStr("6048af620731d1bc6c38cfd85df2d5b194af25acc3b109d2fc572b6460219922")},
		{Height: 2400000, Hash: newHashFromStr("ed5569553f6966c7b6b75e78734a7499eb8b3ef619230262c374ef25e903328a")},
		{Height: 2500000, Hash: newHashFromStr("239a6e786001523e8ba4485b4fc1665e911a5b9b093260d194bb9c9b394d244a")},
		{Height: 2600000, Hash: newHashFromStr("e4cd35742ba676ec33782cd022c66f96db2b93001de073c5fa07f76831517f65")},
		{Height: 2700000, Hash: newHashFromStr("b6f5ad74dc4d013f9d2729341cd559d2141352aa25566b7f0951553487455a1c")},
		{Height: 2800000, Hash: newHashFromStr("99b08a2bb086a355af00782dffc2a94f56e225426426b3748063e4c703444a45")},
		{Height: 2900000, Hash: newHashFromStr("9888c9f637e68fae64a3e6381872394124351611e65dabd77997652c33c3d2f9")},
		{Height: 3000000, Hash: newHashFromStr("c9b034e634cb78f16385ba5cd166f91a5b448af84d6b20c0a924bc2f4409effd")},
		{Height: 3100000, Hash: newHashFromStr("db8d0280c6fc40196cf214e094fe87930fb6cd19d430fede2012539648ffb881")},
		{Height: 3200000, Hash: newHashFromStr("5dcc21cebabb27e1e4376f6018169bfbeda345ace052e4e4f8be363150016858")},
		{Height: 3300000, Hash: newHashFromStr("0bceb218f7d8ed33dfba60a4c6866798367a6d76883ff01cc3194b8f68fa261e")},
		{Height: 3400000, Hash: newHashFromStr("00000000000002656bb99809312027f25a4436ca6d4161b379046a6291566b12")},
		{Height: 3500000, Hash: newHashFromStr("bece76f2a3f53637e2ea84837a45a6ffdc0c86372ab4701c3146094f65832c80")},
		{Height: 3600000, Hash: newHashFromStr("520931769d7ec1fcade7ae754943d844b27adb7c09e7dece6ac40320bf135c09")},
		{Height: 3700000, Hash: newHashFromStr("9e6aa6a32b00f8919fb06bade175ec7e78a1d6cd3cac66095d2f6a2e7cfd7adb")},
		{Height: 3800000, Hash: newHashFromStr("45dc73e94f35017464da16e7a703e5559f7f2e79585a7cbfb048bcaa56eaab67")},
		{Height: 3900000, Hash: newHashFromStr("7f63bbf0296697374e046392c60388af7bb3b05861eba66bfa226a8c4e666375")},
		{Height: 4000000, Hash: newHashFromStr("000000000000009d41478ed798aa84f059430efe0b493c2eedd6a17a6afde1cc")},
		{Height: 4100000, Hash: newHashFromStr("4a4a01554dca8f252dca0dbc166e2f16e22108d3d76afc1ba4718c8b048da355")},
		{Height: 4200000, Hash: newHashFromStr("140ea726d95fdaf39f9a4920ba2346b4165c24a7839c78b7cd9728678c9f42d9")},
		{Height: 4300000, Hash: newHashFromStr("6076bbd15d330f1acaf70a6f3e6c74c71c53616fe83835527910859df68fe302")},
		{Height: 4400000, Hash: newHashFromStr("aeb94a76714dc6577077d02bd048f05f0b1f6d52336674830d692be52f345b0f")},
		{Height: 4500000, Hash: newHashFromStr("000000000000006027d9f6aec51709c4b9e8bc4a0dba1e881847c30e3bf427f4")},
		{Height: 4600000, Hash: newHashFromStr("c17a42929bf68ca3a0fef128b2d2097b32f8bc64025cd0e00494e97ecc031389")},
		{Height: 4700000, Hash: newHashFromStr("b208ab3e67343b03e0ed0b86c1847137edf52255abc0c502dfd6dae730d87bc7")},
		{Height: 4800000, Hash: newHashFromStr("baac90bad8948257c3fcd11986e83ef58d21b5cb9b542bfd97fa07f74cbd4f07")},
		{Height: 4900000, Hash: newHashFromStr("dbdd7a7c0e97e7c9f863a3388e503853572251b67cc5f5a7d019029901f5d910")},
		{Height: 5000000, Hash: newHashFromStr("1dd2fdf6416343688eed463a7bc70b298a4f872e941e36f85cda0915d6488e25")},
		{Height: 5100000, Hash: newHashFromStr("000000000000001fa2a2aa5ea1486bb95aa94359db9051f2a58f7ebba7581648")},
		{Height: 5200000, Hash: newHashFromStr("e219eaa86c4b6d8446f746ed50299829c864100b43c430070639f051922ef98a")},
		{Height: 5300000, Hash: newHashFromStr("8603ca6b0a3e586cb9bb4d3d8c55659a41e2200228ed5f2994a3b648cbdba63b")},
		{Height: 5400000, Hash: newHashFromStr("8925d8f745f1e0bae5381e4a998c2b77e27faacdac6e8c39ab5e9a1ad4283775")},
		{Height: 5500000, Hash: newHashFromStr("e77792fd9f4cfaf3987078d6f70755e13f9c47a8cf854180405598c59daca031")},
		{Height: 5600000, Hash: newHashFromStr("a9cc9e925cf3b65185b255fa00911fb5008ccc4125ca3cd7fda8028d8cdefd62")},
		{Height: 5700000, Hash: newHashFromStr("832ccb9f8ddb429dfedf8fa560eb3f5b1de4731630172d3ce18fe7504beaa15c")},
		{Height: 5800000, Hash: newHashFromStr("80687f64052f45c64cf62e92615732f80fc18b1ee1b47c699a192141699676f0")},
		{Height: 5900000, Hash: newHashFromStr("b35ef100b062ecaff14f1a82f31548ca1da18290648ce4125737c99b4741b41f")},
		{Height: 6000000, Hash: newHashFromStr("6495a84f8f83981a435a6cbf9e6dd4bf0f38618c8325213ca6ef6add40c0ddd8")},
		{Height: 6100000, Hash: newHashFromStr("000000000000000140aad8cbf12752aae2fbc8d835ce1a75b658a45be134456a")},
		{Height: 6200000, Hash: newHashFromStr("f3dac63c3b0b83f8dc5776312ffce138612b4339df4738990858a876c290f7a2")},
		{Height: 6300000, Hash: newHashFromStr("6980190f2e7b5a2a76f1a825570ac229bdcbf819ee244bff2f70f78b52823aca")},
		{Height: 6400000, Hash: newHashFromStr("ec40a9adc7f44de66f04c3c8abb26f58f5714960f1f8506904e621cd9666896c")},
		{Height: 6500000, Hash: newHashFromStr("b168b7f70cbfd2e5fea07da55d9fa90dc7c65599ceb2700efe04ee6c45692e52")},
		{Height: 6600000, Hash: newHashFromStr("fb0fdab30d2fe283d92149e959a14f7cb63eaddf6077797de8922e34e4a5233e")},
		{Height: 6700000, Hash: newHashFromStr("34f0e2229e32ebdc0f4b1642f4737bfab072410906557107003e2b773ac915c0")},
		{Height: 6800000, Hash: newHashFromStr("9e931416e834b8b3da0f9a4bc7e7650683de6cb9d2220f5d716f3eadba785751")},
		{Height: 6900000, Hash: newHashFromStr("31c160cbae37afcacd351597c162661951806861f4b7afa3b2a5050dcec1f82c")},
		{Height: 7000000, Hash: newHashFromStr("03c6664b250c3e3b688f5779ce791384b35acaa38c4461f0458a4674bd762f63")},
		{Height: 7100000, Hash: newHashFromStr("f7485dd054e59567e5ceb8dbe52ff3aef4f8745e5e08a52979fac01748b8572a")},
		{Height: 7200000, Hash: newHashFromStr("8ddd0ef75cb09f84b2d78a128425276e4aa3e52155a19c076c1e697f792cc9ab")},
		{Height: 7300000, Hash: newHashFromStr("7dd0e9e988f475db6c9542895629e3ce7a51febc3dca08c0e663672661ffa50b")},
		{Height: 7400000, Hash: newHashFromStr("00000000000000002376f94b6c7cc3e1c8b76c20f57d9df152b456096793fcb6")},
		{Height: 7500000, Hash: newHashFromStr("bd10f5ebedd8d09b352d9e95c96af62ae0a0a7e3f698c6c82c151bc770c6a831")},
		{Height: 7600000, Hash: newHashFromStr("0faf84b6789cc109e36c2d546904897df9adfa76e45a7d9b4e6fca0985905a4d")},
		{Height: 7700000, Hash: newHashFromStr("13447b2e8d708b83114f00a717f8c64d9a1ab508e8121dd14dee95747776b6d8")},
		{Height: 7800000, Hash: newHashFromStr("88e4802635239b3b87ee26e660e68de46d29aa13515e337f056e612030ef02e8")},
		{Height: 7900000, Hash: newHashFromStr("757958c8be093dcb4c4b225305e27849a9bb01d64258793b400339a570a7f100")},
		{Height: 8000000, Hash: newHashFromStr("1af919cb004bb05c369a862cb5ded70aaa123d0eac2432ceec859f6f42880660")},
		{Height: 8100000, Hash: newHashFromStr("00000000000000059722b0f86997bdad4e0384740eb7620b1071d48528400b49")},
		{Height: 8200000, Hash: newHashFromStr("c048b60f03f23188806956dd8d368c94749fc6d8fe2c5751d63a0301ccfcff3b")},
		{Height: 8300000, Hash: newHashFromStr("c61c84b023ad4b5b16a6e365fa2d632368986fffc1cf9777420cd5d73807c715")},
		{Height: 8400000, Hash: newHashFromStr("0000000000000000126a30d5179f6c11674cc861c038dd79eef3146dcd416911")},
		{Height: 8500000, Hash: newHashFromStr("98d5d78c95c762778d6b3e62dfab1c8a212287628ee43aff337ee45ef3ec250b")},
		{Height: 8600000, Hash: newHashFromStr("24a1b9c95d84cc6a2e4ef1282f4e29bc1eb797c6a3788e9a5c601783b0bdcf77")},
		{Height: 8700000, Hash: newHashFromStr("ca0d3e0572902706a1ec13225fe6597ca6db73a6d04da736526bf1c059f04f51")},
		{Height: 8800000, Hash: newHashFromStr("647fa4a4fcb981300019f2808096e31e7d46bc3463b51ddb11cec3e908252c52")},
		{Height: 8900000, Hash: newHashFromStr("ef459b7767eae0588ca1544a36968ec07259cb048599c8414d5310b77ec43898")},
		{Height: 9000000, Hash: newHashFromStr("942b62f60ae25478d6ee41ec498daefc306cf6f93ff500435a12aa6fe3750220")},
		{Height: 9100000, Hash: newHashFromStr("0000000000000002a481e71184e9922cb0b23b3c0aa38276d406cd19b371d6d0")},
		{Height: 9200000, Hash: newHashFromStr("ff44d27cfaa33cb449bcac144d4a920144fece73ff635ac465788194df7dc8b3")},
		{Height: 9300000, Hash: newHashFromStr("a0f214790be2a6c51401fabe7166536f27f396993394c0522a7ba076b94abe59")},
		{Height: 9400000, Hash: newHashFromStr("231815644a71807ce4eab21bd010eb4933f15c55f67956825dfe37ce95af02f7")},
		{Height: 9500000, Hash: newHashFromStr("5b0351361414e520e9132ba6c5c4926d6f9ee55c41b77fffce3a16ea15d4a1be")},
		{Height: 9600000, Hash: newHashFromStr("00000000000000002b2ff20c6168d24e7722a3ec9e36b69f2f365d4ec5fc51dc")},
		{Height: 9700000, Hash: newHashFromStr("86a47c07e9b931339e2c2006d1be14cbb43fd5fb446c5cf9d3838d6c2ac44ee0")},
		{Height: 9800000, Hash: newHashFromStr("0e52a46580927e519789868f0f033c8cc7c2e67a5878ab3bc79d7d87e68d4c9d")},
		{Height: 9900000, Hash: newHashFromStr("00000000000000008c8789cf584524433adc107a4c5e2747216a739e7b767f69")},
		{Height: 10000000, Hash: newHashFromStr("9e382e2ae1909a4f20c40f38bc7b9f5d0222d5f92ed8be9c04238209c88d55b7")},
		{Height: 10100000, Hash: newHashFromStr("a79d5e1dbf2d4a72791c8941c7bc946198d09ab8a1f82bdcfa7eaa1c741baaac")},
		{Height: 10200000, Hash: newHashFromStr("e14c037f86661f6a7f886e8b658ed7d612ea94ab7202076f14fed723ca8abc12")},
		{Height: 10300000, Hash: newHashFromStr("cf879988b014d1ad51493ed4392273dead9ff950db7cf9d528f302441a6baa4f")},
		{Height: 10400000, Hash: newHashFromStr("ffd5406c454fc3c7483fe5f3fa124e3a0d0a25570719c0fd8cfd39bee4f78687")},
		{Height: 10500000, Hash: newHashFromStr("2e7c794d328830995cded5b6a37a238159ea8f9553665d4e77eb4f13d10794b3")},
		{Height: 10600000, Hash: newHashFromStr("28a97f6a1b2bfeab7a60c265e38957c6f2c7e3db4dc61cb83e681b23f3544075")},
		{Height: 10700000, Hash: newHashFromStr("43c7ce1278e9fc1706bec51029af1a379f8a68e32b1d9eef4562ffd24c42a9a8")},
		{Height: 10800000, Hash: newHashFromStr("13347801d7d0fdb856b150a766e06bb2d4115f1b279f1225e66da0da95aaac91")},
		{Height: 10900000, Hash: newHashFromStr("6ab645a1d1827ead16d1b6ff084af1949bc293815b82e68c79aeaaf164faf156")},
		{Height: 11000000, Hash: newHashFromStr("0f4ad10ae49b504246c0175f6cbab9b0f91b6568a88931e6341a83a731701054")},
		{Height: 11100000, Hash: newHashFromStr("000000000000000142cc4591c0c1b31d95781b947067ea5dd4ca448b63d29bbf")},
		{Height: 11200000, Hash: newHashFromStr("8392bf656466bcbe5b409b9878aa2628a43e1657fccc4ad4900b97df33b71b94")},
	},

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 1916, // 95% of MinerConfirmationWindow
	MinerConfirmationWindow:       2016, //
	Deployments: [DefinedDeployments]chaincfg.ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  27,
			StartTime:  1199145601, // January 1, 2008 UTC
			ExpireTime: 1230767999, // December 31, 2008 UTC
		},
		DeploymentCSV: {
			BitNumber:  12,
			StartTime:  1489997089, // March 24th, 2017
			ExpireTime: 1521891345, // March 24th, 2018
		},
		DeploymentSegwit: {
			BitNumber:  13,
			StartTime:  1490355345, // March 24th, 2017
			ExpireTime: 1521891345, // March 24th, 2018
		},
	},

	// Mempool parameters
	RelayNonStdTxs: false,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "dgb", // always bc for main net

	// Address encoding magics
	PubKeyHashAddrID:        0x1e, // starts with 1
	ScriptHashAddrID:        0x32, // starts with 3
	PrivateKeyID:            0x80, // starts with 5 (uncompressed) or K (compressed)
	WitnessPubKeyHashAddrID: 0x06, // starts with p2
	WitnessScriptHashAddrID: 0x0A, // starts with 7Xh

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 0x14,
}

var RegressionNetParams = chaincfg.Params{
	Name: "regtest",

	// DigiByte has 0xdab5bffa as RegTest (same as Bitcoin's RegTest).
	// Setting it to an arbitrary value (leet_hex(digibyte)), so that we can
	// register the regtest network.
	// DigiByte Core Developers will change this soon.
	Net:         0xd191841e,
	DefaultPort: "18444",
	DNSSeeds:    []chaincfg.DNSSeed{
		// None
	},

	// Chain parameters
	GenesisBlock:             &genesisBlock,
	GenesisHash:              &genesisHash,
	PowLimit:                 new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne),
	PowLimitBits:             0x1d00ffff,
	BIP0034Height:            4394880, // add8ca420f557f62377ec2be6e6f47b96cf2e68160d58aeb7b73433de834cca0
	BIP0065Height:            4394880, // add8ca420f557f62377ec2be6e6f47b96cf2e68160d58aeb7b73433de834cca0
	BIP0066Height:            4394880, // add8ca420f557f62377ec2be6e6f47b96cf2e68160d58aeb7b73433de834cca0
	CoinbaseMaturity:         100,
	SubsidyReductionInterval: 210000,
	TargetTimespan:           12 * time.Hour / 5, // 2.4 hours
	TargetTimePerBlock:       time.Second * 60,   // 60 seconds
	RetargetAdjustmentFactor: 4,                  // 25% less, 400% more
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0,
	GenerateSupported:        false,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []chaincfg.Checkpoint{},

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 1916, // 95% of MinerConfirmationWindow
	MinerConfirmationWindow:       2016, //
	Deployments: [DefinedDeployments]chaincfg.ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  27,
			StartTime:  1199145601, // January 1, 2008 UTC
			ExpireTime: 1230767999, // December 31, 2008 UTC
		},
		DeploymentCSV: {
			BitNumber:  12,
			StartTime:  1489997089, // March 24th, 2017
			ExpireTime: 1521891345, // March 24th, 2018
		},
		DeploymentSegwit: {
			BitNumber:  13,
			StartTime:  1490355345, // March 24th, 2017
			ExpireTime: 1521891345, // March 24th, 2018
		},

		// These Deployments can't not be used, because the struct only allows
		// three Deployments.

		// DeploymentVersionBits: {
		// 	BitNumber: 14,
		// 	StartTime: 1521891345,
		// 	ExpireTime: 1489997089,
		// },

		// DeploymentVersionReserveAlgos: {
		// 	BitNumber: 12,
		// 	StartTime: 1574208000,
		// 	ExpireTime: 1542672000,
		// },

		// DeploymentOdo: {
		// 	BitNumber: 6,
		// 	StartTime: 1588291200,
		// 	ExpireTime: 1556668800,
		// },

		// DeploymentEquihash: {
		// 	BitNumber: 3,
		// 	StartTime: 1521891345,
		// 	ExpireTime: 1489997089,
		// },

		// DeploymentEthash: {
		// 	BitNumber: 4,
		// 	StartTime: 1521891345,
		// 	ExpireTime: 1489997089,
		// },
	},

	// Mempool parameters
	RelayNonStdTxs: false,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "dgbrt", // always bc for main net

	// Address encoding magics
	PubKeyHashAddrID:        0x7e, // starts with 1
	ScriptHashAddrID:        0x8c, // starts with 3
	PrivateKeyID:            0xfe, // starts with 5 (uncompressed) or K (compressed)
	WitnessPubKeyHashAddrID: 0x06, // starts with p2
	WitnessScriptHashAddrID: 0x0A, // starts with 7Xh

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with xpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 0x14,
}
