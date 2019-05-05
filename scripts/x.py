#!/usr/bin/env python
# coding: utf-8

from pykit import mathtoy
from pykit.mathtoy import Polynomial

want_range_size = 15
order = 3

ys = [0, 16, 32, 48, 64, 79, 95, 111, 126, 142, 158, 174, 190, 206, 222, 236,
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
      1330, 1332, 1334, 1336, 1338, 1340, 1342, 1344, 1346, 1348, 1350, 1352]
xs = list(range(len(ys)))

def eval_maxmindiff(poly, xs, ys):
    maxdiff = 0
    mindiff = 0
    for x, y in zip(xs, ys):

        v = Polynomial.evaluate(poly, x)
        diff = y - v
        if maxdiff < diff:
            maxdiff = diff
        if mindiff > diff:
            mindiff = diff

    return maxdiff, mindiff

def resample_points(poly, xs, ys):
    newys = []
    for x, y in zip(xs, ys):
        v = Polynomial.evaluate(poly, x)
        v = round(v)

        diff = int(y -v)
        newys.append(diff)

    return newys

def find_curv(xs, ys, rng):

    l, r = 0, len(xs) + 1

    while l < r - 1:
        mid = (l+r)/2
        xx, yy = xs[:mid], ys[:mid]
        poly = Polynomial.fit(xx, yy, order)
        print "try l, mid, r:", l, mid, r
        print "fit poly:", poly

        maxdiff, mindiff = eval_maxmindiff(poly, xx, yy)
        curr_rng = maxdiff - mindiff
        print "diffs:", maxdiff, mindiff, "range size:", curr_rng

        if curr_rng <= rng:
            l = mid
        else:
            r = mid

    poly = Polynomial.fit(xs[:l], ys[:l], order)
    maxdiff, mindiff = eval_maxmindiff(poly, xs[:l], ys[:l])
    poly[0] += mindiff

    # print "max diff:", curr_rng, rng, curr_rng < rng
    return poly, l


def find_polys(xs, ys):
    polys = []
    points = []
    l = 0
    while len(xs) > 0:
        print
        print "start==="
        print
        poly, n = find_curv(xs, ys, want_range_size)

        print "found a curve:"
        print "    curve:", poly
        print "    range:", l, l +n

        polys.append(poly)
        resampled = resample_points(poly, xs[:n], ys[:n])
        points.append([xs[:n], resampled])
        xs = xs[n:]
        ys = ys[n:]

        l+=n

    return polys, points

rng = [0, xs[-1]]
polys, points = find_polys(xs, ys)
for poly, xy in zip(polys, points):
    print poly, xy[1]

# polys = [polys[4]]
syms = 'abcdefghijk'
polys = [[x, syms[i]] for i, x in enumerate(polys)]

for l in Polynomial.plot(polys, rng,
                         rangey=[0, 1500],
                         height=40,
                         # points=zip(xs, ys)
):
    print l

print len(xs)

raise


poly[0] += mindiff

print poly
print maxdiff - mindiff

for x, y in enumerate(ys):

    v = mathtoy.Polynomial.evaluate(poly, x)
    diff = y - v
    print int(diff), 

print
