import { Game } from './game.js';

const CONTENT = {
    CREATE_GAME: Symbol("CREATE_GAME"),
    IN_QUEUE: Symbol("IN_QUEUE"),
    IN_GAME: Symbol("IN_GAME"),
}

let createGameHTML = `
    <div id="create-game" class="main-container" style="background-color: lightslategray">
        <button id="find-match">Find Match!</button>
    </div>
`

let inQueueHTML = `
    <div id="in-queue" class="main-container" style="background-color: lightsteelblue">
        <button id="find-match">Waiting for opponents..</button>
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

class main {
    constructor() {
        this.currentContent = CONTENT.CREATE_GAME;
        this.mainElem = document.querySelector('main');
    }

    render(content_type) {
        if (content_type == CONTENT.CREATE_GAME) {
            this.mainElem.innerHTML = createGameHTML;

            const findMatchBtn = document.getElementById('find-match');
            if (findMatchBtn) {
                findMatchBtn.addEventListener('click', () => {
                    this.render(CONTENT.IN_QUEUE);
                });
            }

        } else if (content_type == CONTENT.IN_QUEUE) {
            this.mainElem.innerHTML = inQueueHTML;
            this.initiateSocket();
        } else if (content_type == CONTENT.IN_GAME) {
            this.mainElem.innerHTML = inGameHTML;
        } else {
            throw new Error("Invalid content");
        }
    }

    initiateSocket() {
        const socket = new WebSocket("ws://localhost:8080/ws");

        socket.onopen = () => {
            console.log("Connected");
        };

        socket.onmessage = (event) => {
            console.log("Received from server:", event.data);
            this.handleGameCreation(socket, event);
        };

        socket.onerror = (err) => {
            console.error("WebSocket error:", err);
        };
    }

    handleGameCreation(socket, event) {
        this.render(CONTENT.IN_GAME);
        const data = JSON.parse(event.data);
        const playerBoard = document.getElementById('player-board');
        const opponentBoard = document.getElementById('opponents-board');
        const [p1RevealedDiv, p2RevealedDiv] = document.querySelectorAll('#score-board .score');

        this.game = new Game(socket, data.width, data.height, playerBoard, opponentBoard, p1RevealedDiv, p2RevealedDiv);
    }
}

let content = new main();
content.render(CONTENT.CREATE_GAME);
