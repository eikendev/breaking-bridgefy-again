const PATH_LIST = 'java.util.List';

const PATH_BLEENTITY = 'com.bridgefy.sdk.framework.entities.BleEntity';
const PATH_FORWARDPACKET = 'com.bridgefy.sdk.framework.entities.ForwardPacket';
const PATH_FORWARDTRANSACTION = 'com.bridgefy.sdk.framework.entities.ForwardTransaction';
const PATH_SESSION = 'com.bridgefy.sdk.framework.controller.Session';

const ENTITY_TYPE_MESH = 3;

const RECEIVER_TYPE_CONTACT = 0;
const RECEIVER_TYPE_BROADCAST = 1;

function printMessages(bleEntity) {
	const JList = Java.use(PATH_LIST);
	const JForwardPacket = Java.use(PATH_FORWARDPACKET);
	const JForwardTransaction = Java.use(PATH_FORWARDTRANSACTION);

	if (bleEntity.getEt() == ENTITY_TYPE_MESH) {
		const transaction = Java.cast(bleEntity.getCt(), JForwardTransaction);
		const packets = Java.cast(transaction.getMesh(), JList);

		for (var i = 0; i < packets.size(); i++) {
			const packet = Java.cast(packets.get(i), JForwardPacket);

			if (packet.getReceiver_type() == RECEIVER_TYPE_CONTACT) {
				const sender = packet.getSender();
				const receiver = packet.getReceiver();

				console.log('[+] ' + sender + ' -> ' + receiver);
			} else if (packet.getReceiver_type() == RECEIVER_TYPE_BROADCAST) {
				const sender = packet.getSender();
				const receiver = 'all';

				console.log('[+] ' + sender + ' -> ' + receiver);
			}
		}
	}
}

function run() {
	const JSession = Java.use(PATH_SESSION);

	const Fa = JSession.a.overload(PATH_BLEENTITY);
	Fa.implementation = function (bleEntity) {
		printMessages(bleEntity);

		return this.a(bleEntity);
	}
};

setImmediate(function () { Java.perform(run) });
