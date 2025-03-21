/*
 * Emby Server REST API
 *
 * Explore the Emby Server API
 *
 */
package embyclient

// ColorFormats : Enum listing color formats.      The enum member names are matching those that are used by ffmpeg. (execute 'ffmpeg \\-pix\\_fmts' to list them) Exception: Items that are starting with a number are prefixed with an underscore here. In ffmpeg code these are prefixed with 'AV\\_PIX\\_FMT\\_' and all\\-caps.
type ColorFormats string

// List of ColorFormats
const (
	UNKNOWN_ColorFormats          ColorFormats = "Unknown"
	YUV420P_ColorFormats          ColorFormats = "yuv420p"
	YUYV422_ColorFormats          ColorFormats = "yuyv422"
	RGB24_ColorFormats            ColorFormats = "rgb24"
	BGR24_ColorFormats            ColorFormats = "bgr24"
	YUV422P_ColorFormats          ColorFormats = "yuv422p"
	YUV444P_ColorFormats          ColorFormats = "yuv444p"
	YUV410P_ColorFormats          ColorFormats = "yuv410p"
	YUV411P_ColorFormats          ColorFormats = "yuv411p"
	GRAY_ColorFormats             ColorFormats = "gray"
	MONOW_ColorFormats            ColorFormats = "monow"
	MONOB_ColorFormats            ColorFormats = "monob"
	PAL8_ColorFormats             ColorFormats = "pal8"
	YUVJ420P_ColorFormats         ColorFormats = "yuvj420p"
	YUVJ422P_ColorFormats         ColorFormats = "yuvj422p"
	YUVJ444P_ColorFormats         ColorFormats = "yuvj444p"
	UYVY422_ColorFormats          ColorFormats = "uyvy422"
	UYYVYY411_ColorFormats        ColorFormats = "uyyvyy411"
	BGR8_ColorFormats             ColorFormats = "bgr8"
	BGR4_ColorFormats             ColorFormats = "bgr4"
	BGR4_BYTE_ColorFormats        ColorFormats = "bgr4_byte"
	RGB8_ColorFormats             ColorFormats = "rgb8"
	RGB4_ColorFormats             ColorFormats = "rgb4"
	RGB4_BYTE_ColorFormats        ColorFormats = "rgb4_byte"
	NV12_ColorFormats             ColorFormats = "nv12"
	NV21_ColorFormats             ColorFormats = "nv21"
	ARGB_ColorFormats             ColorFormats = "argb"
	RGBA_ColorFormats             ColorFormats = "rgba"
	ABGR_ColorFormats             ColorFormats = "abgr"
	BGRA_ColorFormats             ColorFormats = "bgra"
	GRAY16_ColorFormats           ColorFormats = "gray16"
	YUV440P_ColorFormats          ColorFormats = "yuv440p"
	YUVJ440P_ColorFormats         ColorFormats = "yuvj440p"
	YUVA420P_ColorFormats         ColorFormats = "yuva420p"
	RGB48_ColorFormats            ColorFormats = "rgb48"
	RGB565_ColorFormats           ColorFormats = "rgb565"
	RGB555_ColorFormats           ColorFormats = "rgb555"
	BGR565_ColorFormats           ColorFormats = "bgr565"
	BGR555_ColorFormats           ColorFormats = "bgr555"
	VAAPI_MOCO_ColorFormats       ColorFormats = "vaapi_moco"
	VAAPI_IDCT_ColorFormats       ColorFormats = "vaapi_idct"
	VAAPI_VLD_ColorFormats        ColorFormats = "vaapi_vld"
	YUV420P16_ColorFormats        ColorFormats = "yuv420p16"
	YUV422P16_ColorFormats        ColorFormats = "yuv422p16"
	YUV444P16_ColorFormats        ColorFormats = "yuv444p16"
	DXVA2_VLD_ColorFormats        ColorFormats = "dxva2_vld"
	RGB444_ColorFormats           ColorFormats = "rgb444"
	BGR444_ColorFormats           ColorFormats = "bgr444"
	YA8_ColorFormats              ColorFormats = "ya8"
	BGR48_ColorFormats            ColorFormats = "bgr48"
	YUV420P9_ColorFormats         ColorFormats = "yuv420p9"
	YUV420P10_ColorFormats        ColorFormats = "yuv420p10"
	YUV422P10_ColorFormats        ColorFormats = "yuv422p10"
	YUV444P9_ColorFormats         ColorFormats = "yuv444p9"
	YUV444P10_ColorFormats        ColorFormats = "yuv444p10"
	YUV422P9_ColorFormats         ColorFormats = "yuv422p9"
	GBRP_ColorFormats             ColorFormats = "gbrp"
	GBRP9_ColorFormats            ColorFormats = "gbrp9"
	GBRP10_ColorFormats           ColorFormats = "gbrp10"
	GBRP16_ColorFormats           ColorFormats = "gbrp16"
	YUVA422P_ColorFormats         ColorFormats = "yuva422p"
	YUVA444P_ColorFormats         ColorFormats = "yuva444p"
	YUVA420P9_ColorFormats        ColorFormats = "yuva420p9"
	YUVA422P9_ColorFormats        ColorFormats = "yuva422p9"
	YUVA444P9_ColorFormats        ColorFormats = "yuva444p9"
	YUVA420P10_ColorFormats       ColorFormats = "yuva420p10"
	YUVA422P10_ColorFormats       ColorFormats = "yuva422p10"
	YUVA444P10_ColorFormats       ColorFormats = "yuva444p10"
	YUVA420P16_ColorFormats       ColorFormats = "yuva420p16"
	YUVA422P16_ColorFormats       ColorFormats = "yuva422p16"
	YUVA444P16_ColorFormats       ColorFormats = "yuva444p16"
	VDPAU_ColorFormats            ColorFormats = "vdpau"
	XYZ12_ColorFormats            ColorFormats = "xyz12"
	NV16_ColorFormats             ColorFormats = "nv16"
	NV20_ColorFormats             ColorFormats = "nv20"
	RGBA64_ColorFormats           ColorFormats = "rgba64"
	BGRA64_ColorFormats           ColorFormats = "bgra64"
	YVYU422_ColorFormats          ColorFormats = "yvyu422"
	YA16_ColorFormats             ColorFormats = "ya16"
	GBRAP_ColorFormats            ColorFormats = "gbrap"
	GBRAP16_ColorFormats          ColorFormats = "gbrap16"
	QSV_ColorFormats              ColorFormats = "qsv"
	MMAL_ColorFormats             ColorFormats = "mmal"
	D3D11VA_VLD_ColorFormats      ColorFormats = "d3d11va_vld"
	CUDA_ColorFormats             ColorFormats = "cuda"
	RGB__ColorFormats             ColorFormats = "_0rgb"
	RGB0_ColorFormats             ColorFormats = "rgb0"
	BGR__ColorFormats             ColorFormats = "_0bgr"
	BGR0_ColorFormats             ColorFormats = "bgr0"
	YUV420P12_ColorFormats        ColorFormats = "yuv420p12"
	YUV420P14_ColorFormats        ColorFormats = "yuv420p14"
	YUV422P12_ColorFormats        ColorFormats = "yuv422p12"
	YUV422P14_ColorFormats        ColorFormats = "yuv422p14"
	YUV444P12_ColorFormats        ColorFormats = "yuv444p12"
	YUV444P14_ColorFormats        ColorFormats = "yuv444p14"
	GBRP12_ColorFormats           ColorFormats = "gbrp12"
	GBRP14_ColorFormats           ColorFormats = "gbrp14"
	YUVJ411P_ColorFormats         ColorFormats = "yuvj411p"
	BAYER_BGGR8_ColorFormats      ColorFormats = "bayer_bggr8"
	BAYER_RGGB8_ColorFormats      ColorFormats = "bayer_rggb8"
	BAYER_GBRG8_ColorFormats      ColorFormats = "bayer_gbrg8"
	BAYER_GRBG8_ColorFormats      ColorFormats = "bayer_grbg8"
	BAYER_BGGR16_ColorFormats     ColorFormats = "bayer_bggr16"
	BAYER_RGGB16_ColorFormats     ColorFormats = "bayer_rggb16"
	BAYER_GBRG16_ColorFormats     ColorFormats = "bayer_gbrg16"
	BAYER_GRBG16_ColorFormats     ColorFormats = "bayer_grbg16"
	XVMC_ColorFormats             ColorFormats = "xvmc"
	YUV440P10_ColorFormats        ColorFormats = "yuv440p10"
	YUV440P12_ColorFormats        ColorFormats = "yuv440p12"
	AYUV64_ColorFormats           ColorFormats = "ayuv64"
	VIDEOTOOLBOX_VLD_ColorFormats ColorFormats = "videotoolbox_vld"
	P010_ColorFormats             ColorFormats = "p010"
	GBRAP12_ColorFormats          ColorFormats = "gbrap12"
	GBRAP10_ColorFormats          ColorFormats = "gbrap10"
	MEDIACODEC_ColorFormats       ColorFormats = "mediacodec"
	GRAY12_ColorFormats           ColorFormats = "gray12"
	GRAY10_ColorFormats           ColorFormats = "gray10"
	GRAY14_ColorFormats           ColorFormats = "gray14"
	P016_ColorFormats             ColorFormats = "p016"
	D3D11_ColorFormats            ColorFormats = "d3d11"
	GRAY9_ColorFormats            ColorFormats = "gray9"
	GBRPF32_ColorFormats          ColorFormats = "gbrpf32"
	GBRAPF32_ColorFormats         ColorFormats = "gbrapf32"
	DRM_PRIME_ColorFormats        ColorFormats = "drm_prime"
	OPENCL_ColorFormats           ColorFormats = "opencl"
	GRAYF32_ColorFormats          ColorFormats = "grayf32"
	YUVA422P12_ColorFormats       ColorFormats = "yuva422p12"
	YUVA444P12_ColorFormats       ColorFormats = "yuva444p12"
	NV24_ColorFormats             ColorFormats = "nv24"
	NV42_ColorFormats             ColorFormats = "nv42"
)
