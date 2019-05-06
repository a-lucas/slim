package polyfit_test

import (
	"fmt"
	"testing"

	"github.com/openacid/slim/benchhelper"
	"github.com/openacid/slim/polyfit"
)

var nums []int32 = []int32{0, 16, 32, 48, 64, 79, 95, 111, 126, 142, 158, 174, 190, 206, 222, 236,
	252, 268, 275, 278, 281, 283, 285, 289, 296, 301, 304, 307, 311, 313, 318,
	321, 325, 328, 335, 339, 344, 348, 353, 357, 360, 364, 369, 372, 377, 383,
	387, 393, 399, 404, 407, 410, 415, 418, 420, 422, 426, 430, 434, 439, 444,
	446, 448, 451, 456, 459, 462, 465, 470, 473, 479, 482, 488, 490, 494, 500,
	506, 509, 513, 519, 521, 528, 530, 534, 537, 540, 544, 546, 551, 556, 560,
	566, 568, 572, 574, 576, 580, 585, 588, 592, 594, 600, 603, 606, 608, 610,
	614, 620, 623, 628, 630, 632, 638, 644, 647, 653, 658, 660, 662, 665, 670,
	672, 676, 681, 683, 687, 689, 691, 693, 695, 697, 703, 706, 710, 715, 719,
	722, 726, 731, 735, 737, 741, 748, 750, 753, 757, 763, 766, 768, 775, 777,
	782, 785, 791, 795, 798, 800, 806, 811, 815, 818, 821, 824, 829, 832, 836,
	838, 842, 846, 850, 855, 860, 865, 870, 875, 878, 882, 886, 890, 895, 900,
	906, 910, 913, 916, 921, 925, 929, 932, 937, 940, 942, 944, 946, 952, 954,
	956, 958, 962, 966, 968, 971, 975, 979, 983, 987, 989, 994, 997, 1000,
	1003, 1008, 1014, 1017, 1024, 1028, 1032, 1034, 1036, 1040, 1044, 1048,
	1050, 1052, 1056, 1058, 1062, 1065, 1068, 1072, 1078, 1083, 1089, 1091,
	1094, 1097, 1101, 1104, 1106, 1110, 1115, 1117, 1119, 1121, 1126, 1129,
	1131, 1134, 1136, 1138, 1141, 1143, 1145, 1147, 1149, 1151, 1153, 1155,
	1157, 1159, 1161, 1164, 1166, 1168, 1170, 1172, 1174, 1176, 1178, 1180,
	1182, 1184, 1186, 1189, 1191, 1193, 1195, 1197, 1199, 1201, 1203, 1205,
	1208, 1210, 1212, 1214, 1217, 1219, 1221, 1223, 1225, 1227, 1229, 1231,
	1233, 1235, 1237, 1239, 1241, 1243, 1245, 1247, 1249, 1251, 1253, 1255,
	1257, 1259, 1261, 1263, 1265, 1268, 1270, 1272, 1274, 1276, 1278, 1280,
	1282, 1284, 1286, 1288, 1290, 1292, 1294, 1296, 1298, 1300, 1302, 1304,
	1306, 1308, 1310, 1312, 1314, 1316, 1318, 1320, 1322, 1324, 1326, 1328,
	1330, 1332, 1334, 1336, 1338, 1340, 1342, 1344, 1346, 1348, 1350, 1352}

func TestFoo(t *testing.T) {

	xs := []float64{1, 2, 3, 4}
	ys := []float64{6, 5, 7, 10}

	// coef = Polynomial.fit(xs, ys, degree=1)
	// self.assertEqual([3.5, 1.4], coef)

	// coef = Polynomial.fit(xs, ys, degree=2)
	// self.assertEqual([8.5, -3.6, 1], coef)
	polyfit.PolyFit(xs, ys, 1)
	polyfit.PolyFit(xs, ys, 2)

	for _, width := range []uint{1, 2, 4, 8} {

		rst := polyfit.Resample32To4(nums, width)
		// for i, rg := range rst.Regions {
		//     fmt.Println("region: ", i, polyfit.PolyStr(rg.Poly),
		//         polyfit.SizeOf(rg))
		// }

		for i, n := range nums {
			r := rst.Get(int32(i))
			if r != n {
				t.Fatalf("i=%d expect: %v; but: %v", i, n, r)
			}
		}

		fmt.Println(width, benchhelper.SizeOf(rst.Datas))
		fmt.Println(width, benchhelper.SizeOf(rst), benchhelper.SizeOf(nums))
	}

	// // Initialize two matrices, a and ia.
	// a := mat.NewDense(2, 2, []float64{
	//     4, 0,
	//     0, 4,
	// })
	// var ia mat.Dense

	// // Take the inverse of a and place the result in ia.
	// ia.Inverse(a)

	// // Print the result using the formatter.
	// fa := mat.Formatted(&ia, mat.Prefix("     "), mat.Squeeze())
	// fmt.Printf("ia = %.2g\n\n", fa)

	// // Confirm that A * A^-1 = I
	// var r mat.Dense
	// r.Mul(a, &ia)
	// fr := mat.Formatted(&r, mat.Prefix("    "), mat.Squeeze())
	// fmt.Printf("r = %v\n\n", fr)

	// // The Inverse operation, however, is numerically unstable,
	// // and should typically be avoided.
	// // For example, a common need is to find x = A^-1 * b.
	// // In this case, the SolveVec method of VecDense
	// // (if b is a Vector) or Solve method of Dense (if b is a
	// // matrix) should used instead of computing the Inverse of A.
	// b := mat.NewDense(2, 2, []float64{
	//     2, 0,
	//     0, 2,
	// })
	// var x mat.Dense
	// x.Solve(a, b)

	// // Print the result using the formatter.
	// fx := mat.Formatted(&x, mat.Prefix("    "), mat.Squeeze())
	// fmt.Printf("x = %v", fx)
}

var Output int32

func BenchmarkGet(b *testing.B) {

	width := uint(4)
	l := len(nums)

	rst := polyfit.Resample32To4(nums, width)

	b.ResetTimer()

	s := int32(0)
	for i := 0; i < b.N; i++ {
		s += rst.Get(int32(i % l))
	}

	Output = s
}
