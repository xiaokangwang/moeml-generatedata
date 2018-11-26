package challangegen

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"path/filepath"
	"time"
	"github.com/disintegration/imaging"
)

type Generator struct {
	BackgroundDir string
	ForegroundDir string
	ForegroundItemSum int
	OutDir string
}
func (g* Generator)GetRandomFileInDir(dir string)string{
	mat,err:=filepath.Glob(dir+"/*.png")
	if err!=nil {
		panic(err)
	}
	add:=rand.Intn(len(mat))
	return mat[add]
}
func (g* Generator)GetRandomBackground()*image.NRGBA {
	result := g.GetRandomFileInDir(g.BackgroundDir)
	img, err := imaging.Open(result)
	if err != nil {
		panic(err)
	}
	return imaging.Clone(img)
}
func (g* Generator)GetRandomForeground()*image.NRGBA{

	result := g.GetRandomFileInDir(g.ForegroundDir)
	img, err := imaging.Open(result)
	if err != nil {
		panic(err)
	}
	return imaging.Clone(img)
}
func (g* Generator)ComposeAll(ite int){
	rand.Seed(time.Now().UnixNano())

	b:=g.GetRandomBackground()
	a:=imaging.New(b.Bounds().Size().X,b.Bounds().Size().Y,color.RGBA{255,255,255,255})

	for i := 0; i<=g.ForegroundItemSum;i++ {

		f:=g.GetRandomForeground()
		ft:=g.RandTransform(f)
		b,a=g.Compose(b,ft,a)

	}
	imaging.Save(b,fmt.Sprintf("%v/%v.png",g.OutDir,ite))
	imaging.Save(a,fmt.Sprintf("%v/%v-a.png",g.OutDir,ite))
}

func (g* Generator) RandTransform(inputimg *image.NRGBA)*image.NRGBA{
	inputimgo := inputimg

	//FlipH pass
	if rand.Intn(1)==1 {
		inputimgo = imaging.FlipH(inputimgo)
	}


	//Scale pass
	ratio:=rand.Float64()*0.6
	inputimgo = imaging.Resize(inputimgo,
		int(ratio*float64(inputimgo.Bounds().Size().X)),
		int(ratio*float64(inputimgo.Bounds().Size().Y)),
		imaging.NearestNeighbor)

	//Rotate pass
	inputimgo = imaging.Rotate(inputimgo,rand.Float64()*360,color.NRGBA{0,0,0,0})

	return inputimgo

}

func (g* Generator)Compose(back,fore,a *image.NRGBA)(*image.NRGBA,*image.NRGBA){
	out:=imaging.Clone(back)
	outAlpha:=imaging.Clone(a)
	bsize:=out.Bounds()
	fsize:=fore.Bounds()

	//find a random location
	liberityX:=bsize.Size().X
	liberityY:=bsize.Size().Y

	moveX:=rand.Intn(liberityX) - fsize.Size().X/2
	moveY:=rand.Intn(liberityY) - fsize.Size().Y/2

	for cX:=0;cX<bsize.Size().X;cX++ {
		for cY:=0;cY<bsize.Size().Y;cY++ {
			if cX > moveX && cX < moveX + fsize.Size().X {
				if cY > moveY && cY < moveY + fsize.Size().Y {
					orig:=out.At(cX,cY)
					XSrc:=cX-moveX
					YSrc:=cY-moveY
					overlay:=fore.At(XSrc,YSrc)
					origR,origG,origB,origA:=orig.RGBA()
					overlayR,overlayG,overlayB,overlayA:=overlay.RGBA()
					origRF:=float64(origR)/float64(0xffff)
					origGF:=float64(origG)/float64(0xffff)
					origBF:=float64(origB)/float64(0xffff)
					origAF:=float64(origA)/float64(0xffff)

					_=origAF

					overlayRF:=float64(overlayR)/float64(0xffff)
					overlayGF:=float64(overlayG)/float64(0xffff)
					overlayBF:=float64(overlayB)/float64(0xffff)
					overlayAF:=float64(overlayA)/float64(0xffff)

					outRF:=origRF*(1-overlayAF)+overlayRF*overlayAF
					outGF:=origGF*(1-overlayAF)+overlayGF*overlayAF
					outBF:=origBF*(1-overlayAF)+overlayBF*overlayAF

					outAF:=1.0

					outAFA:=1.0-overlayAF


					outR:=uint8(outRF*float64(0xff))
					outG:=uint8(outGF*float64(0xff))
					outB:=uint8(outBF*float64(0xff))
					outA:=uint8(outAF*float64(0xff))


					newRGBA:=color.RGBA{outR,outG,outB,outA}

					_=newRGBA

					out.Set(cX,cY,newRGBA)

					AC :=outAlpha.At(cX,cY)
					acR,acG,acB,acA:=AC.RGBA()


					ACRF:=float64(acR)/float64(0xffff)
					ACGF:=float64(acG)/float64(0xffff)
					ACBF:=float64(acB)/float64(0xffff)
					ACAF:=float64(acA)/float64(0xffff)

					_=ACAF
					ACRF=ACRF*outAFA
					ACGF=ACGF*outAFA
					ACBF=ACBF*outAFA
					ACAF=1.0

					outACR:=uint8(ACRF*float64(0xff))
					outACG:=uint8(ACGF*float64(0xff))
					outACB:=uint8(ACBF*float64(0xff))
					outACA:=uint8(ACAF*float64(0xff))

					newAC:=color.RGBA{outACR,outACG,outACB,outACA}

					outAlpha.Set(cX,cY,newAC)
				}

			}
		}
	}

	return out,outAlpha
}
