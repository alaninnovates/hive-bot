package hiveplugin

import (
	"io"
	"strconv"
	"strings"

	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/common/loaders"
	"alaninnovates.com/hive-bot/hiveplugin/hive"
	"github.com/disgoorg/disgo/discord"
	"github.com/fogleman/gg"
)

func GetRangeNumbers(rangeStr string) []int {
	var nums []int
	for _, s := range strings.Split(rangeStr, ",") {
		if strings.Contains(s, "-") {
			l, _ := strconv.ParseInt(strings.Split(s, "-")[0], 10, 64)
			r, _ := strconv.ParseInt(strings.Split(s, "-")[1], 10, 64)
			for i := l; i <= r; i++ {
				nums = append(nums, int(i))
			}
		} else {
			n, _ := strconv.ParseInt(s, 10, 64)
			nums = append(nums, int(n))
		}
	}
	return nums
}

var InvalidSlotsMessage = discord.MessageCreate{
	Content: "Invalid slot range. Slots must be whole numbers between 1 and 50.",
}

func ValidateRange(rangeStr string, min int, max int) bool {
	for _, num := range GetRangeNumbers(rangeStr) {
		if num < min || num > max {
			return false
		}
	}
	return true
}

func RenderHiveImage(h *hive.Hive, showHiveNumbers bool, slotsOnTop bool, skipHiveNumbers []int, background string) *io.PipeReader {
	dc := gg.NewContext(410, 950)
	hive.DrawHive(h, dc, showHiveNumbers, slotsOnTop, skipHiveNumbers)
	img := dc.Image()
	bg, err := gg.LoadImage(loaders.GetHiveBackgroundImagePath(background))
	if err != nil {
		bg, _ = gg.LoadImage(loaders.GetHiveBackgroundImagePath("default"))
	}
	hiveImage := gg.NewContextForImage(bg)
	hiveImage.DrawImageAnchored(img, hiveImage.Width()/2, hiveImage.Height()/2, 0.5, 0.5)
	return common.ImageToPipe(hiveImage.Image())
}
