const CELL_SIZE = "24px";
const CELLS = {
    "9": "images/button_mine.jpg",
    "0": "images/button_0.jpg",
    "1": "images/button_1.jpg",
    "2": "images/button_2.jpg",
    "3": "images/button_3.jpg",
    "4": "images/button_4.jpg",
    "5": "images/button_5.jpg",
    "6": "images/button_6.jpg",
    "7": "images/button_7.jpg",
    "8": "images/button_8.jpg",
}

export class Game {
    constructor(socket, width, height, p1Board, p2Board) {
        this.socket = socket;
        this.playerBoard = p1Board;
        this.opponentBoard = p2Board;
        this.width = width;
        this.height = height;

        this.renderMinesweeperBoard(this.playerBoard, true);
        this.renderMinesweeperBoard(this.opponentBoard, false);

        this.socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            this.processEvent(data);
        };
    }

    renderMinesweeperBoard(board, interactable) {
        board.style.display = 'grid';
        board.style.gridTemplateColumns = `repeat(${this.width}, ${CELL_SIZE})`;
        board.style.gridTemplateRows = `repeat(${this.height}, ${CELL_SIZE})`;

        board.innerHTML = '';
        for (let i = 0; i < this.height; i++) {
            for (let j = 0; j < this.width; j++) {
                const cell = document.createElement('img');
                cell.src = 'images/button_unrevealed.jpg';
                cell.draggable = false;

                cell.classList.add('player-cell');
                cell.row = i;
                cell.col = j;
                cell.flagged = false;

                cell.addEventListener('contextmenu', (e) => e.preventDefault());

                if (interactable) {
                    cell.addEventListener('mousedown', (e) => {
                        e.preventDefault();
                        if (e.button === 2) {
                            this.placeFlag(cell);
                        } else if (e.button === 0) {
                            this.interact(cell);
                        }
                    });
                }
                board.appendChild(cell);
            }
        }
    }

    interact(cell) {
        const message = { x: cell.col, y: cell.row };
        console.log("Sending move to server:", message);
        this.socket.send(JSON.stringify(message));
    }

    processEvent(data) {
        if (data.positions) {
            for (const pos of data.positions) {
                if (data.yourMove) {
                    this.revealCellPlayer(pos.y, pos.x, pos.value);
                } else {
                    this.revealCellOpponent(pos.y, pos.x, pos.value);
                }
            }
        }

        if (data.hitMine && data.yourMove) {
            alert("You hit a mine!");
        } else if (data.won && data.yourMove) {
            alert("You won!");
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
    }

    revealCellOpponent(row, col, value) {
        const cell = this.getCellOpponent(row, col);
        cell.src = CELLS[value];
    }

    placeFlag(cell) {
        if (!cell.flagged) {
            cell.src = 'images/button_flag.jpg';
            cell.flagged = true;
        } else {
            cell.src = 'images/button_unrevealed.jpg';
            cell.flagged = false;
        }
    }
}
