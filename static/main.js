import { Game } from './game.js';

const CONTENT = {
    CREATE_GAME: Symbol("CREATE_GAME"),
    IN_QUEUE: Symbol("IN_QUEUE"),
    IN_GAME: Symbol("IN_GAME"),
    LADDER: Symbol("LADDER"),
    ABOUT: Symbol("ABOUT"),
    REGISTER: Symbol("REGISTER")
}

let createGameHTML = `
    <div id="create-game" class="main-container" style="background-color: lightslategray">
        <button id="find-match-button">Find Match!</button>
    </div>
`

let inQueueHTML = `
    <div id="in-queue" class="main-container" style="background-color: lightsteelblue">
        <button id="in-queue-button">Waiting for opponents..</button>
        <div id="in-queue-share-msg">Looks like nobody is playing at the moment. Send this link to a friend!</p>
    </div>
`;

let inGameHTML = `
    <div id="in-game" class="main-container" style="background-color: lightyellow">
        <div id="player-container" style="background-color: #ebf0b1">
            <div id="player-board" style="background-color: #c3cad6"></div>
        </div>

        <div id="board-separator" style="background-color: black"></div>

        <div id="opponent-container" style="background-color: #cacf91">
            <div id="opponents-board" style="background-color: #d0ddf2"></div>
        </div>

        </div>
    </div>
`;

let ladderHTML = `
    <div id="ladder" class="main-container" style="background-color: lightyellow">
        Ladder unavailable
    </div>
`

let aboutHTML = `
    <div id="about" class="main-container" style="background-color: lightyellow">
        <p>Created for Czech Technical University,</p>
        <p>Faculty of Electrical Engineering,</p>
        <p>Vývoj klientských aplikací v Javascriptu</p>
        <p>(B0B39KAJ)</p>
        <p>2025, Marek Jagoš</p>
    </div>
`

let registerHTML = `
    <div id="register" class="main-container" style="background-color: lightyellow">
        Registration unavailable
    </div>
`

class main {
    constructor() {
        this.mainElem = document.querySelector('main');
        this.queueTimeout = null;
        this.gameStartSound = new Audio('/sounds/game-start.mp3');
    }

    render(content_type) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN && content_type != CONTENT.IN_GAME) {
            this.socket.close();
            this.socket = null;
            console.log("Socket closed");
        }

        if (content_type == CONTENT.CREATE_GAME) {
            this.mainElem.innerHTML = createGameHTML;

            const findMatchBtn = document.getElementById('find-match-button');
            findMatchBtn.addEventListener('click', () => {
                this.render(CONTENT.IN_QUEUE);
            });
            this.currentContent = CONTENT.CREATE_GAME;

        } else if (content_type == CONTENT.IN_QUEUE) {
            this.mainElem.innerHTML = inQueueHTML;
            this.initiateSocket();

            const inQueueButton = document.getElementById('in-queue-button');
            inQueueButton.addEventListener('click', () => {
                this.render(CONTENT.CREATE_GAME);
            });

            if (this.queueTimeout) {
                clearTimeout(this.queueTimeout);
            }

            this.queueTimeout = setTimeout(() => {
                const msgElem = document.getElementById('in-queue-share-msg');
                if (this.currentContent === CONTENT.IN_QUEUE) {
                    msgElem.style.display = 'block';
                }
            }, 5000);

            this.currentContent = CONTENT.IN_QUEUE;

        } else if (content_type == CONTENT.IN_GAME) {
            this.mainElem.innerHTML = inGameHTML;
            this.currentContent = CONTENT.IN_GAME;

        } else if (content_type == CONTENT.LADDER) {
            this.mainElem.innerHTML = ladderHTML;
            this.currentContent = CONTENT.LADDER;

        } else if (content_type == CONTENT.ABOUT) {
            this.mainElem.innerHTML = aboutHTML;
            this.currentContent = CONTENT.ABOUT;

        } else if (content_type == CONTENT.REGISTER) {
            this.mainElem.innerHTML = registerHTML;
            this.currentContent = CONTENT.REGISTER;

        } else {
            throw new Error("Invalid content");
        }
    }

    initiateSocket() {
        const protocol = window.location.protocol === "https:" ? "wss" : "ws";
        const host = window.location.host;
        const socketUrl = `${protocol}://${host}/ws`;
        this.socket = new WebSocket(socketUrl);

        this.socket.onopen = () => {
            console.log("Connected");
        };

        this.socket.onmessage = (event) => {
            console.log("Received from server:", event.data);
            this.handleGameCreation(this.socket, event);
        };

        this.socket.onerror = (err) => {
            console.error("WebSocket error:", err);
        };

        this.socket.onclose = () => {
            console.log("WebSocket closed");
        };
    }

    handleGameCreation(socket, event) {
        this.render(CONTENT.IN_GAME);
        const data = JSON.parse(event.data);
        console.log(data);

        const playerBoard = document.getElementById('player-board');
        const opponentBoard = document.getElementById('opponents-board');
        const [p1RevealedDiv, p2RevealedDiv] = document.querySelectorAll('#score-board .score');

        this.game = new Game(socket, data.width, data.height, data.startX, data.startY, playerBoard, opponentBoard, p1RevealedDiv, p2RevealedDiv);
        this.gameStartSound.play();
    }
}

let content = new main();
content.render(CONTENT.CREATE_GAME);

const hamburgerBtn = document.getElementById('hamburger-menu');
const aside = document.querySelector('.page section aside');

hamburgerBtn.addEventListener('click', () => {
    const isHidden = getComputedStyle(aside).display === 'none';
    aside.style.display = isHidden ? 'flex' : 'none';
});

window.addEventListener('resize', () => {
    if (window.innerWidth > 768) {
        aside.style.display = 'flex';
    } else {
        aside.style.display = 'none';
    }
})

let loggedIn = false;
let user = "Unknown";
document.querySelectorAll('.menu-item').forEach(item => {
    item.addEventListener('click', (e) => {
        const action = item.getAttribute('data-action');

        switch (action) {
            case 'play':
                content.render(CONTENT.CREATE_GAME);
                break;
            case 'ladder':
                content.render(CONTENT.LADDER);
                break;
            case 'about':
                content.render(CONTENT.ABOUT);
                break;
            case 'register':
                if (!loggedIn) {
                    content.render(CONTENT.REGISTER);
                }
                break;
            default:
                console.log('Unknown action');
        }
    });
});