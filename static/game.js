const CELL_SIZE = "24px";

const MINE = 9;
const FLAG = 10;
const UNREVEALED = 11;
const START = 12;

const CELLS = {
    0: "images/0.jpg",
    1: "images/1.jpg",
    2: "images/2.jpg",
    3: "images/3.jpg",
    4: "images/4.jpg",
    5: "images/5.jpg",
    6: "images/6.jpg",
    7: "images/7.jpg",
    8: "images/8.jpg",
    [MINE]: "images/mine.jpg",
    [FLAG]: "images/flag.jpg",
    [UNREVEALED]: "images/unrevealed.jpg",
    [START]: "images/start.jpg",
}

export class Game {
    constructor(socket, width, height, startX, startY, p1Board, p2Board) {
        this.socket = socket;
        this.playerBoard = p1Board;
        this.opponentBoard = p2Board;
        this.width = width;
        this.height = height;
        this.startX = startX;
        this.startY = startY;
        this.gameConcluded = false;

        this.renderMinesweeperBoard(this.playerBoard, true);
        this.renderMinesweeperBoard(this.opponentBoard, false);

        this.socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            console.log("Received from server: ", data);
            this.processEvent(data);
        };

        this.gameWinSound = new Audio('/sounds/game-win.mp3');
        this.gameLoseSound = new Audio('/sounds/game-lose.mp3');
    }

    renderMinesweeperBoard(board, interactable) {
        board.style.display = 'grid';
        board.style.gridTemplateColumns = `repeat(${this.width}, ${CELL_SIZE})`;
        board.style.gridTemplateRows = `repeat(${this.height}, ${CELL_SIZE})`;

        board.innerHTML = '';
        for (let i = 0; i < this.height; i++) {
            for (let j = 0; j < this.width; j++) {
                const cell = document.createElement('img');
                cell.src = CELLS[UNREVEALED];
                cell.draggable = false;

                cell.classList.add('player-cell');
                cell.row = i;
                cell.col = j;
                cell.flagged = false;

                cell.addEventListener('contextmenu', (e) => e.preventDefault());

                if (interactable) {
                    cell.addEventListener('mousedown', (e) => {
                        e.preventDefault();
                        this.interact(cell, e.button == 2);
                    });
                }
                board.appendChild(cell);
            }
        }

        // Set starting cell
        const cell = this.getCellPlayer(this.startX, this.startY);
        cell.src = CELLS[START];
    }

    interact(cell, flag) {
        const message = { x: cell.col, y: cell.row, flag: flag };
        console.log("Sending move to server:", message);
        this.socket.send(JSON.stringify(message));
    }

    processEvent(data) {
        if (!this.gameConcluded && (data.win || data.lose)) {
            const msgElem = document.getElementById('lose-wait-msg');
            msgElem.style.display = 'none';
            this.gameConcluded = true;

            if (data.win) {
                this.gameWinSound.play();
                const msgElemWin = document.getElementById('win-conclude-msg');
                msgElemWin.style.display = 'block';
                
            } else if (data.lose) {
                this.gameLoseSound.play();
                const msgElemLose = document.getElementById('lose-conclude-msg');
                msgElemLose.style.display = 'block';
            }
        }

        for (const pos of data.cellUpdates) {
            if (data.yourMove) {
                this.revealCellPlayer(pos.y, pos.x, pos.value);
            } else {
                this.revealCellOpponent(pos.y, pos.x, pos.value);
            }
        }
    }

    getCellPlayer(row, col) {
        return this.playerBoard.children[row * this.width + col];
    }

    getCellOpponent(row, col) {
        return this.opponentBoard.children[row * this.width + col];
    }

    revealCellPlayer(row, col, value) {
        const cell = this.getCellPlayer(row, col);
        cell.src = CELLS[value];

        if (value == MINE) {
            if (!this.gameConcluded) {
                this.queueTimeout = setTimeout(() => {
                const msgElem = document.getElementById('lose-wait-msg');
                msgElem.style.display = 'block';
                }, 2000);
            }
        }
    }

    revealCellOpponent(row, col, value) {
        const cell = this.getCellOpponent(row, col);
        cell.src = CELLS[value];
    }
}
