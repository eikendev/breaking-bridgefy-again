const PATH_HEX = 'org.apache.commons.codec.binary.Hex';
const PATH_STRING = 'java.lang.String';
const PATH_BLUETOOTHGATTCHARACTERISTIC = 'android.bluetooth.BluetoothGattCharacteristic';

const PATH_UTILS = 'com.bridgefy.sdk.framework.utils.Utils';
const PATH_UTIL = 'com.bridgefy.sdk.framework.crypto.Util';
const PATH_CHUNKUTILS = 'com.bridgefy.sdk.framework.controller.ChunkUtils';

const DEBUG = false;

function encodeHex(byteArray) {
	let HexClass = Java.use(PATH_HEX);
	let StringClass = Java.use(PATH_STRING);
	const hexChars = HexClass.encodeHex(byteArray);
	return StringClass.$new(hexChars).toString();
}

function run() {
	console.log('[*] Overriding implementations...')

	let JUtils = Java.use(PATH_UTILS);
	let JUtil = Java.use(PATH_UTIL);
	let JChunkUtils = Java.use(PATH_CHUNKUTILS);

	JUtils.fromMessagePacktoEntity.implementation = function (data, clazz) {
		console.log('[+] Util.fromMessagePacktoEntity()');
		return this.fromMessagePacktoEntity(data, clazz);
	};

	JUtil.decrypt.implementation = function (data, key) {
		console.log('[+] Util.decrypt()');
		return this.decrypt(data, key);
	};

	const Fa = JChunkUtils.a.overload('[B');
	Fa.implementation = function (data) {
		console.log('[+] ChunkUtils.a(len = ' + data.length + ') // before add()');
		const stripped = this.a(data);

		if (DEBUG)
			console.log(encodeHex(stripped));

		return stripped;
	};

	const Fa2 = JChunkUtils.a.overload(PATH_BLUETOOTHGATTCHARACTERISTIC);
	Fa2.implementation = function (data) {
		console.log('[+] ChunkUtils.a() // after add()');
		return this.a(data);
	};

	JChunkUtils.decompress.implementation = function (data) {
		console.log('[+] ChunkUtils.decompress() enter');
		const decompressed = this.decompress(data);
		console.log('[+] ChunkUtils.decompress() exit');

		return decompressed;
	};
}

setImmediate(function () { Java.perform(run) });
