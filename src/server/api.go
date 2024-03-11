package server

import (
	"SilentPaymentAppBackend/src/common"
	"SilentPaymentAppBackend/src/db/dblevel"
	"bytes"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

// todo might not need ApiHandler struct if no data is stored within.
//  Will keep for now just in case, so I don't have to refactor twice
type ApiHandler struct{}

type TxRequest struct {
	Data string `form:"data" json:"data" binding:"required"`
}

func (h *ApiHandler) GetBestBlockHeight(c *gin.Context) {
	lastHeader, err := dblevel.FetchHighestBlockHeaderInv()
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could could not retrieve data from database",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"block_height": lastHeader.Height,
	})
}

func (h *ApiHandler) GetCFilterByHeight(c *gin.Context) {
	heightStr := c.Param("blockheight")
	if heightStr == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	height, err := strconv.ParseUint(heightStr, 10, 32)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not parse height",
		})
		return
	}
	headerInv, err := dblevel.FetchByBlockHeightBlockHeaderInv(uint32(height))
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get height mapping from db",
		})
		return
	}

	cFilter, err := dblevel.FetchByBlockHashFilter(headerInv.Hash)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get filter from db",
		})
		return
	}

	data := gin.H{
		"filter_type":  cFilter.FilterType,
		"block_height": height, // saves us a "join" in the query
		"block_header": cFilter.BlockHash,
		"data":         hex.EncodeToString(cFilter.Data),
	}

	c.JSON(200, data)
}

func (h *ApiHandler) GetLightUTXOsByHeight(c *gin.Context) {
	heightStr := c.Param("blockheight")
	if heightStr == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	height, err := strconv.ParseUint(heightStr, 10, 32)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not parse height",
		})
		return
	}
	headerInv, err := dblevel.FetchByBlockHeightBlockHeaderInv(uint32(height))
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get height mapping from db",
		})
		return
	}
	utxos, err := dblevel.FetchByBlockHashUTXOs(headerInv.Hash)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could could not retrieve data from database",
		})
		return
	}
	if utxos != nil {
		c.JSON(200, utxos)
	} else {
		c.JSON(200, []interface{}{})
	}
}

// GetTweakDataByHeight serves tweak data as json array of tweaks (33 byte as hex-formatted)
// todo can be changed to serve with verbosity aka serve with txid or even block data (height, hash)
func (h *ApiHandler) GetTweakDataByHeight(c *gin.Context) {
	heightStr := c.Param("blockheight")
	if heightStr == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	height, err := strconv.ParseUint(heightStr, 10, 32)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not parse height",
		})
		return
	}
	headerInv, err := dblevel.FetchByBlockHeightBlockHeaderInv(uint32(height))
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get height mapping from db",
		})
		return
	}
	tweaks, err := dblevel.FetchByBlockHashTweaks(headerInv.Hash)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could could not retrieve data from database",
		})
		return
	}
	var serveTweakData []string
	for _, tweak := range tweaks {
		serveTweakData = append(serveTweakData, hex.EncodeToString(tweak.Data[:]))
	}

	c.JSON(200, serveTweakData)
}

func (h *ApiHandler) ForwardRawTX(c *gin.Context) {
	var txRequest TxRequest
	if err := c.ShouldBind(&txRequest); err != nil {
		common.ErrorLogger.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}
	err := forwardTxToMemPool(txRequest.Data)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func forwardTxToMemPool(txHex string) error {
	//url := "http://localhost/api/tx"

	resp, err := http.Post(common.MempoolEndpoint, "application/x-www-form-urlencoded", bytes.NewBufferString(txHex))
	if err != nil {
		common.ErrorLogger.Printf("Failed to make request: %s\n", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		common.ErrorLogger.Printf("Failed to read response: %s\n", err)
		return err
	}

	common.DebugLogger.Println("Response:", string(body))
	return nil
}
