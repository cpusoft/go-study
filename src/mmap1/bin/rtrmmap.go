package main

import (
	_ "github.com/cpusoft/goutil/logs"
)

/*
// 测试整体流程：写入 Full + 写入 Incremental → 读取并校验
func main() {
	ctx := context.Background()
	zaplogs.InfoArgs(ctx, "main")
	// ======================
	// 1. 生成测试数据
	// ======================
	testFullList := genTestRtrFull(5)      // 生成 5 条 Full 测试数据
	testIncrMap := genTestRtrIncremental() // 生成 2 个版本增量数据
	zaplogs.InfoArgs(ctx, "main", "len(testFullList)", len(testFullList),
		"len(testIncrMap)", len(testIncrMap))
	// ======================
	// 2. 写入 MMap
	// ======================
	zaplogs.DebugArgs(ctx, "\n=== 开始写入 RTR Full MMap ===")
	err := mmap.WriteRtrFullMmap(ctx, testFullList)
	if err != nil {
		zaplogs.DebugArgs(ctx, "WriteRtrFullMmap 失败", "err", err)
		return
	}

	zaplogs.DebugArgs(ctx, "\n=== 开始写入 RTR Incremental MMap ===")
	err = mmap.WriteRtrIncrementalMmap(ctx, testIncrMap)
	if err != nil {
		zaplogs.DebugArgs(ctx, "WriteRtrIncrementalMmap 失败", "err", err)
		return
	}

	// ======================
	// 3. 读取 MMap
	// ======================
	zaplogs.DebugArgs(ctx, "\n=== 开始读取 RTR Full MMap ===")
	readFullList, err := mmap.ReadRtrFullMmap(ctx)
	if err != nil {
		zaplogs.DebugArgs(ctx, "ReadRtrFullMmap 失败", "err", err)
		return
	}

	zaplogs.DebugArgs(ctx, "\n=== 开始读取 RTR Incremental MMap ===")
	readIncrMap, err := mmap.ReadRtrIncrementalMmap(ctx)
	if err != nil {
		zaplogs.DebugArgs(ctx, "ReadRtrIncrementalMmap 失败", "err", err)
		return
	}

	// ======================
	// 4. 校验数据一致性
	// ======================
	zaplogs.DebugArgs(ctx, "=== 校验 Full 数据：写入 条，读取 条",
		"len(testFullList)", len(testFullList), "len(readFullList)", len(readFullList))
	if len(testFullList) != len(readFullList) {
		zaplogs.DebugArgs(ctx, "Full 数据条数不匹配")
	}
	zaplogs.DebugArgs(ctx, "Full 数据条数校验通过")

	zaplogs.DebugArgs(ctx, "=== 校验 Incremental 数据：写入  个版本，读取  个版本",
		"len(testIncrMap)", len(testIncrMap), "len(readIncrMap)", len(readIncrMap))
	if len(testIncrMap) != len(readIncrMap) {
		zaplogs.DebugArgs(ctx, "Incremental 版本数不匹配")
	}
	zaplogs.DebugArgs(ctx, "Incremental 版本数校验通过")

	// 逐条对比 Full 数据
	for i := range testFullList {
		orig := testFullList[i]
		read := readFullList[i]
		if orig.Asn != read.Asn ||
			orig.Address != read.Address ||
			orig.PrefixLength != read.PrefixLength ||
			orig.MaxLength != read.MaxLength {
			zaplogs.DebugArgs(ctx, "Full 第  条数据不匹配", "i", i)
		}
	}
	zaplogs.DebugArgs(ctx, "Full 逐条数据校验通过")

	zaplogs.DebugArgs(ctx, "\n========================")
	zaplogs.DebugArgs(ctx, "所有 MMap 读写测试全部通过！")
	zaplogs.DebugArgs(ctx, "========================")
}

// 生成测试用 RTR Full 数据
func genTestRtrFull(count int) []LabRpkiRtrFull {
	var list []LabRpkiRtrFull
	for i := 0; i < count; i++ {
		list = append(list, LabRpkiRtrFull{
			Asn:          int64(1000 + i),
			Address:      fmt.Sprintf("203.10.10.%d", i),
			PrefixLength: 24,
			MaxLength:    32,
		})
	}
	return list
}

// 生成测试用 RTR Incremental 数据（2 个版本）
func genTestRtrIncremental() map[uint64][]LabRpkiRtrIncremental {
	incMap := make(map[uint64][]LabRpkiRtrIncremental)

	// 版本 100：2 条
	incMap[100] = []LabRpkiRtrIncremental{
		{
			SerialNumber: 100,
			Asn:          64501,
			Address:      "192.168.1.10",
			PrefixLength: 24,
			MaxLength:    24,
			Style:        "1", // announce
		},
		{
			SerialNumber: 100,
			Asn:          64502,
			Address:      "192.168.2.10",
			PrefixLength: 24,
			MaxLength:    24,
			Style:        "0", // withdraw
		},
	}

	// 版本 200：1 条
	incMap[200] = []LabRpkiRtrIncremental{
		{
			SerialNumber: 200,
			Asn:          64510,
			Address:      "2001:db8::1",
			PrefixLength: 64,
			MaxLength:    64,
			Style:        "1",
		},
	}

	return incMap
}

// lab_rpki_rtr_full
type LabRpkiRtrFull struct {
	Id           uint64 `json:"id" xorm:"id bigint"`
	SerialNumber uint64 `json:"serialNumber" xorm:"serialNumber bigint"`
	Asn          int64  `json:"asn" xorm:"asn bigint"`
	//address: 63.60.00.00
	Address      string `json:"address" xorm:"address varchar(512)"`
	PrefixLength uint64 `json:"prefixLength" xorm:"prefixLength int"`
	MaxLength    uint64 `json:"maxLength" xorm:"maxLength int"`
	//'come from : {souce:sync/slurm/transfer,syncLogId/syncLogFileId/slurmId/slurmFileId/transferLogId}',
	SourceFrom string `json:"sourceFrom" xorm:"sourceFrom json"`
}
type LabRpkiRtrIncremental struct {
	Id           uint64 `json:"id" xorm:"id bigint"`
	SerialNumber uint64 `json:"serialNumber" xorm:"serialNumber bigint"`
	//announce/withdraw, is 1/0 in protocol
	Style string `json:"style" xorm:"style varchar(16)"`
	Asn   int64  `json:"asn" xorm:"asn bigint"`
	//address: 63.60.00.00
	Address      string `json:"address" xorm:"address varchar(512)"`
	PrefixLength uint64 `json:"prefixLength" xorm:"prefixLength int"`
	MaxLength    uint64 `json:"maxLength" xorm:"maxLength int"`
	//'come from : {souce:sync/slurm/transfer,syncLogId/syncLogFileId/slurmId/slurmFileId/transferLogId}',
	SourceFrom string `json:"sourceFrom" xorm:"sourceFrom json"`
}
*/
