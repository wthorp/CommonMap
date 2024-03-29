package rpf

// This file implements conventions on RPF (CIB, CADRG, DTED)
// data series information and file naming conventions

//Type is RPF data types
type Type int

const (
	//CADRG is Compressed ARC Digitized Raster Graphics
	CADRG Type = iota
	//CIB is Controlled Image Base
	CIB
	//CDTED is Compressed Digital Terrain Elevation Data
	CDTED
)

//NitfSeries captures critical information about a RPF data series
type NitfSeries struct {
	SeriesCode, GroupCode, ScaleText, Name string
	Type                                   Type
	Scale                                  float64
}

//DataSeries enumerates all RPF data series
var DataSeries = map[string]NitfSeries{
	"A1": {"A1", "CM", "1:10K", "Combat Charts (1:10K)", CADRG, 10000},
	"A2": {"A2", "CM", "1:25K", "Combat Charts (1:25K)", CADRG, 25000},
	"A3": {"A3", "CM", "1:50K", "Combat Charts (1:50K)", CADRG, 50000},
	"A4": {"A4", "CM", "1:100K", "Combat Charts (1:100K)", CADRG, 100000},
	"AT": {"AT", "ATC", "1:200K", "Series 200 Air Target Chart", CADRG, 200000},
	"C1": {"C1", "CG", "1:10000", "City Graphics", CADRG, 10000},
	"C2": {"C2", "CG", "1:10560", "City Graphics", CADRG, 10560},
	"C3": {"C3", "CG", "1:11000", "City Graphics", CADRG, 11000},
	"C4": {"C4", "CG", "1:11800", "City Graphics", CADRG, 11800},
	"C5": {"C5", "CG", "1:12000", "City Graphics", CADRG, 12000},
	"C6": {"C6", "CG", "1:12500", "City Graphics", CADRG, 12500},
	"C7": {"C7", "CG", "1:12800", "City Graphics", CADRG, 12800},
	"C8": {"C8", "CG", "1:14000", "City Graphics", CADRG, 14000},
	"C9": {"C9", "CG", "1:14700", "City Graphics", CADRG, 14700},
	"CA": {"CA", "CG", "1:15000", "City Graphics", CADRG, 15000},
	"CB": {"CB", "CG", "1:15500", "City Graphics", CADRG, 15500},
	"CC": {"CC", "CG", "1:16000", "City Graphics", CADRG, 16000},
	"CD": {"CD", "CG", "1:16666", "City Graphics", CADRG, 16666},
	"CE": {"CE", "CG", "1:17000", "City Graphics", CADRG, 17000},
	"CF": {"CF", "CG", "1:17500", "City Graphics", CADRG, 17500},
	"CG": {"CG", "CG", "Various", "City Graphics", CADRG, -1},
	"CH": {"CH", "CG", "1:18000", "City Graphics", CADRG, 18000},
	"CJ": {"CJ", "CG", "1:20000", "City Graphics", CADRG, 20000},
	"CK": {"CK", "CG", "1:21000", "City Graphics", CADRG, 21000},
	"CL": {"CL", "CG", "1:21120", "City Graphics", CADRG, 21120},
	"CM": {"CM", "CM", "Various", "Combat Charts", CADRG, -1},
	"CN": {"CN", "CG", "1:22000", "City Graphics", CADRG, 22000},
	"CO": {"CO", "CO", "Various", "Coastal Charts", CADRG, -1},
	"CP": {"CP", "CG", "1:23000", "City Graphics", CADRG, 23000},
	"CQ": {"CQ", "CG", "1:25000", "City Graphics", CADRG, 25000},
	"CR": {"CR", "CG", "1:26000", "City Graphics", CADRG, 26000},
	"CS": {"CS", "CG", "1:35000", "City Graphics", CADRG, 35000},
	"CT": {"CT", "CG", "1:36000", "City Graphics", CADRG, 36000},
	"D1": {"D1", "DTED1", "100m", "Elevation Data from DTED level 1", CDTED, 100.0},
	"D2": {"D2", "DTED2", "30m", "Elevation Data from DTED level 2", CDTED, 30.0},
	"EG": {"EG", "NARC", "1:11M", "North Atlantic Route Chart", CADRG, 11000000},
	"ES": {"ES", "SEC", "1:500K", "VFR Sectional", CADRG, 500000},
	"ET": {"ET", "SEC", "1:250K", "VFR Sectional Inserts", CADRG, 250000},
	"F1": {"F1", "TFC-1", "1:250K", "Transit Flying Chart (TBD #1)", CADRG, 250000},
	"F2": {"F2", "TFC-2", "1:250K", "Transit Flying Chart (TBD #2)", CADRG, 250000},
	"F3": {"F3", "TFC-3", "1:250K", "Transit Flying Chart (TBD #3)", CADRG, 250000},
	"F4": {"F4", "TFC-4", "1:250K", "Transit Flying Chart (TBD #4)", CADRG, 250000},
	"F5": {"F5", "TFC-5", "1:250K", "Transit Flying Chart (TBD #5)", CADRG, 250000},
	"GN": {"GN", "GNC", "1:5M", "Global Navigation Chart", CADRG, 5000000},
	"HA": {"HA", "HA", "Various", "Harbor and Approach Charts", CADRG, -1},
	"I1": {"I1", "CIB10", "10m", "Imagery, 10 meter resolution", CIB, 10.0},
	"I2": {"I2", "CIB5", "5m", "Imagery, 5 meter resolution", CIB, 5.0},
	"I3": {"I3", "CIB2", "2m", "Imagery, 2 meter resolution", CIB, 2.0},
	"I4": {"I4", "CIB1", "1m", "Imagery, 1 meter resolution", CIB, 1.0},
	"I5": {"I5", "CIB.5", ".5m", "Imagery, .5 (half) meter resolution", CIB, 0.5},
	"IV": {"IV", "", "Various > 10m", "Imagery, greater than 10 meter resolution", CIB, -1},
	"JA": {"JA", "JOG-A", "1:250K", "Joint Operation Graphic - Air", CADRG, 250000},
	"JG": {"JG", "JOG", "1:250K", "Joint Operation Graphic", CADRG, 250000},
	"JN": {"JN", "JNC", "1:2M", "Jet Navigation Chart", CADRG, 2000000},
	"JO": {"JO", "OPG", "1:250K", "Operational Planning Graphic", CADRG, 250000},
	"JR": {"JR", "JOG-R", "1:250K", "Joint Operation Graphic - Radar", CADRG, 250000},
	"K1": {"K1", "ICM", "1:8K", "Image City Maps", CADRG, 8000},
	"K2": {"K2", "ICM", "1:10K", "Image City Maps", CADRG, 10000},
	"K3": {"K3", "ICM", "1:10560", "Image City Maps", CADRG, 10560},
	"K7": {"K7", "ICM", "1:12500", "Image City Maps", CADRG, 12500},
	"K8": {"K8", "ICM", "1:12800", "Image City Maps", CADRG, 12800},
	"KB": {"KB", "ICM", "1:15K", "Image City Maps", CADRG, 15000},
	"KE": {"KE", "ICM", "1:16666", "Image City Maps", CADRG, 16666},
	"KM": {"KM", "ICM", "1:21120", "Image City Maps", CADRG, 21120},
	"KR": {"KR", "ICM", "1:25K", "Image City Maps", CADRG, 25000},
	"KS": {"KS", "ICM", "1:26K", "Image City Maps", CADRG, 26000},
	"KU": {"KU", "ICM", "1:36K", "Image City Maps", CADRG, 36000},
	"L1": {"L1", "LFC-1", "1:500K", "Low Flying Chart (TBD #1)", CADRG, 500000},
	"L2": {"L2", "LFC-2", "1:500K", "Low Flying Chart (TBD #2)", CADRG, 500000},
	"L3": {"L3", "LFC-3", "1:500K", "Low Flying Chart (TBD #3)", CADRG, 500000},
	"L4": {"L4", "LFC-4", "1:500K", "Low Flying Chart (TBD #4)", CADRG, 500000},
	"L5": {"L5", "LFC-5", "1:500K", "Low Flying Chart (TBD #5)", CADRG, 500000},
	"LF": {"LF", "LFC-FR (Day)", "1:500K", "Low Flying Chart (Day) - Host Nation", CADRG, 500000},
	"LN": {"LN", "LN (Night)", "1:500K", "Low Flying Chart (Night) - Host Nation", CADRG, 500000},
	"M1": {"M1", "MIM", "Various", "Military Installation Maps (TBD #1)", CADRG, -1},
	"M2": {"M2", "MIM", "Various", "Military Installation Maps (TBD #2)", CADRG, -1},
	"MH": {"MH", "MIM", "1:25K", "Military Installation Maps", CADRG, 25000},
	"MI": {"MI", "MIM", "1:50K", "Military Installation Maps", CADRG, 50000},
	"MJ": {"MJ", "MIM", "1:100K", "Military Installation Maps", CADRG, 100000},
	"MM": {"MM", "", "Various", "(Miscellaneous Maps & Charts)", CADRG, -1},
	"OA": {"OA", "OPAREA", "Various", "Naval Range Operation Area Chart", CADRG, -1},
	"OH": {"OH", "VHRC", "1:1M", "VFR Helicopter Route Chart", CADRG, 1000000},
	"ON": {"ON", "ONC", "1:1M", "Operational Navigation Chart", CADRG, 1000000},
	"OW": {"OW", "WAC", "1:1M", "High Flying Chart - Host Nation", CADRG, 1000000},
	"P1": {"P1", "", "1:25K", "Special Military Map - Overlay", CADRG, 25000},
	"P2": {"P2", "", "1:25K", "Special Military Purpose", CADRG, 25000},
	"P3": {"P3", "", "1:25K", "Special Military Purpose", CADRG, 25000},
	"P4": {"P4", "", "1:25K", "Special Military Purpose", CADRG, 25000},
	"P5": {"P5", "", "1:50K", "Special Military Map - Overlay", CADRG, 50000},
	"P6": {"P6", "", "1:50K", "Special Military Purpose", CADRG, 50000},
	"P7": {"P7", "", "1:50K", "Special Military Purpose", CADRG, 50000},
	"P8": {"P8", "", "1:50K", "Special Military Purpose", CADRG, 50000},
	"P9": {"P9", "", "1:100K", "Special Military Map - Overlay", CADRG, 100000},
	"PA": {"PA", "", "1:100K", "Special Military Purpose", CADRG, 100000},
	"PB": {"PB", "", "1:100K", "Special Military Purpose", CADRG, 100000},
	"PC": {"PC", "", "1:100K", "Special Military Purpose", CADRG, 100000},
	"PD": {"PD", "", "1:250K", "Special Military Map - Overlay", CADRG, 250000},
	"PE": {"PE", "", "1:250K", "Special Military Purpose", CADRG, 250000},
	"PF": {"PF", "", "1:250K", "Special Military Purpose", CADRG, 250000},
	"PG": {"PG", "", "1:250K", "Special Military Purpose", CADRG, 250000},
	"PH": {"PH", "", "1:500K", "Special Military Map - Overlay", CADRG, 500000},
	"PI": {"PI", "", "1:500K", "Special Military Purpose", CADRG, 500000},
	"PJ": {"PJ", "", "1:500K", "Special Military Purpose", CADRG, 500000},
	"PK": {"PK", "", "1:500K", "Special Military Purpose", CADRG, 500000},
	"PL": {"PL", "", "1:1M", "Special Military Map - Overlay", CADRG, 1000000},
	"PM": {"PM", "", "1:1M", "Special Military Purpose", CADRG, 1000000},
	"PN": {"PN", "", "1:1M", "Special Military Purpose", CADRG, 1000000},
	"PO": {"PO", "", "1:1M", "Special Military Purpose", CADRG, 1000000},
	"PP": {"PP", "", "1:2M", "Special Military Map - Overlay", CADRG, 2000000},
	"PQ": {"PQ", "", "1:2M", "Special Military Purpose", CADRG, 2000000},
	"PR": {"PR", "", "1:2M", "Special Military Purpose", CADRG, 2000000},
	"PS": {"PS", "", "1:5M", "Special Military Map - Overlay", CADRG, 5000000},
	"PT": {"PT", "", "1:5M", "Special Military Purpose", CADRG, 5000000},
	"PU": {"PU", "", "1:5M", "Special Military Purpose", CADRG, 5000000},
	"PV": {"PV", "", "1:5M", "Special Military Purpose", CADRG, 5000000},
	"R1": {"R1", "", "1:50K", "Range Charts", CADRG, 50000},
	"R2": {"R2", "", "1:100K", "Range Charts", CADRG, 100000},
	"R3": {"R3", "", "1:250K", "Range Charts", CADRG, 250000},
	"R4": {"R4", "", "1:500K", "Range Charts", CADRG, 500000},
	"R5": {"R5", "", "1:1M", "Range Charts", CADRG, 1000000},
	"RC": {"RC", "RGS-100", "1:100K", "Russian General Staff Maps", CADRG, 100000},
	"RL": {"RL", "RGS-50", "1:50K", "Russian General Staff Maps", CADRG, 50000},
	"RR": {"RR", "RGS-200", "1:200K", "Russian General Staff Maps", CADRG, 200000},
	"RV": {"RV", "Riverine", "1:50K", "Riverine Map 1:50,000 scale", CADRG, 50000},
	"TC": {"TC", "TLM 100", "1:100K", "Topographic Line Map 1:100,000 scale", CADRG, 100000},
	"TF": {"TF", "TFC (Day)", "1:250K", "Transit Flying Chart (Day)", CADRG, 250000},
	"TL": {"TL", "TLM50", "1:50K", "Topographic Line Map", CADRG, 50000},
	"TN": {"TN", "TFC (Night)", "1:250K", "Transit Flying Chart (Night) - Host Nation", CADRG, 250000},
	"TP": {"TP", "TPC", "1:500K", "Tactical Pilotage Chart", CADRG, 500000},
	"TQ": {"TQ", "TLM24", "1:24K", "Topographic Line Map 1:24,000 scale", CADRG, 24000},
	"TR": {"TR", "TLM200", "1:200K", "Topographic Line Map 1:200,000 scale", CADRG, 200000},
	"TT": {"TT", "TLM25", "1:25K", "Topographic Line Map 1:25,000 scale", CADRG, 25000},
	"UL": {"UL", "TLM50 - Other", "1:50K", "Topographic Line Map (other 1:50,000 scale)", CADRG, 50000},
	"V1": {"V1", "Inset HRC", "1:50", "Helicopter Route Chart Inset", CADRG, 50},
	"V2": {"V2", "Inset HRC", "1:62500", "Helicopter Route Chart Inset", CADRG, 62500},
	"V3": {"V3", "Inset HRC", "1:90K", "Helicopter Route Chart Inset", CADRG, 90000},
	"V4": {"V4", "Inset HRC", "1:250K", "Helicopter Route Chart Inset", CADRG, 250000},
	"VH": {"VH", "HRC", "1:125K", "Helicopter Route Chart", CADRG, 125000},
	"VN": {"VN", "VNC", "1:500K", "Visual Navigation Charts", CADRG, 500000},
	"VT": {"VT", "VTAC", "1:250K", "VFR Terminal Area Chart", CADRG, 250000},
	"WA": {"WA", "", "1:250K", "IFR Enroute Low", CADRG, 250000},
	"WB": {"WB", "", "1:500K", "IFR Enroute Low", CADRG, 500000},
	"WC": {"WC", "", "1:750K", "IFR Enroute Low", CADRG, 750000},
	"WD": {"WD", "", "1:1M", "IFR Enroute Low", CADRG, 1100000},
	"WE": {"WE", "", "1:1.5M", "IFR Enroute Low", CADRG, 1500000},
	"WF": {"WF", "", "1:2M", "IFR Enroute Low", CADRG, 2000000},
	"WG": {"WG", "", "1:2.5M", "IFR Enroute Low", CADRG, 2500000},
	"WH": {"WH", "", "1:3M", "IFR Enroute Low", CADRG, 3000000},
	"WI": {"WI", "", "1:3.5M", "IFR Enroute Low", CADRG, 3500000},
	"WK": {"WK", "", "1:4M", "IFR Enroute Low", CADRG, 4000000},
	"XD": {"XD", "", "1:1M", "IFR Enroute High", CADRG, 1000000},
	"XE": {"XE", "", "1:1.5M", "IFR Enroute High", CADRG, 1500000},
	"XF": {"XF", "", "1:2M", "IFR Enroute High", CADRG, 2000000},
	"XG": {"XG", "", "1:2.5M", "IFR Enroute High", CADRG, 2500000},
	"XH": {"XH", "", "1:3M", "IFR Enroute High", CADRG, 3000000},
	"XI": {"XI", "", "1:3.5M", "IFR Enroute High", CADRG, 3500000},
	"XJ": {"XJ", "", "1:4M", "IFR Enroute High", CADRG, 4000000},
	"XK": {"XK", "", "1:4.5M", "IFR Enroute High", CADRG, 4500000},
	"Y9": {"Y9", "", "1:16.5M", "IFR Enroute Area", CADRG, 16500000},
	"YA": {"YA", "", "1:250K", "IFR Enroute Area", CADRG, 250000},
	"YB": {"YB", "", "1:500K", "IFR Enroute Area", CADRG, 500000},
	"YC": {"YC", "", "1:750K", "IFR Enroute Area", CADRG, 750000},
	"YD": {"YD", "", "1:1M", "IFR Enroute Area", CADRG, 1000000},
	"YE": {"YE", "", "1:1.5M", "IFR Enroute Area", CADRG, 1500000},
	"YF": {"YF", "", "1:2M", "IFR Enroute Area", CADRG, 2000000},
	"YI": {"YI", "", "1:3.5M", "IFR Enroute Area", CADRG, 3500000},
	"YJ": {"YJ", "", "1:4M", "IFR Enroute Area", CADRG, 4000000},
	"YZ": {"YZ", "", "1:12M", "IFR Enroute Area", CADRG, 12000000},
	"ZA": {"ZA", "", "1:250K", "IFR Enroute High/Low", CADRG, 250000},
	"ZB": {"ZB", "", "1:500K", "IFR Enroute High/Low", CADRG, 500000},
	"ZC": {"ZC", "", "1:750K", "IFR Enroute High/Low", CADRG, 750000},
	"ZD": {"ZD", "", "1:1M", "IFR Enroute High/Low", CADRG, 1000000},
	"ZE": {"ZE", "", "1:1.5M", "IFR Enroute High/Low", CADRG, 1500000},
	"ZF": {"ZF", "", "1:2M", "IFR Enroute High/Low", CADRG, 2000000},
	"ZG": {"ZG", "", "1:2.5M", "IFR Enroute High/Low", CADRG, 2500000},
	"ZH": {"ZH", "", "1:3M", "IFR Enroute High/Low", CADRG, 3000000},
	"ZI": {"ZI", "", "1:3.5M", "IFR Enroute High/Low", CADRG, 3500000},
	"ZJ": {"ZJ", "", "1:4M", "IFR Enroute High/Low", CADRG, 4000000},
	"ZK": {"ZK", "", "1:4.5M", "IFR Enroute High/Low", CADRG, 4500000},
	"ZT": {"ZT", "", "1:9M", "IFR Enroute High/Low", CADRG, 9000000},
	"ZV": {"ZV", "", "1:10M", "IFR Enroute High/Low", CADRG, 10000000},
	"ZZ": {"ZZ", "", "1:12M", "IFR Enroute High/Low", CADRG, 12000000},
}
