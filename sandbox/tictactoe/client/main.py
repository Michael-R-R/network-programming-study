#!/usr/bin/python

import sys, json

try:
    import PyQt6.QtCore
    import PyQt6.QtWidgets
    import PyQt6.QtGui
    import PyQt6.QtNetwork
except Exception as e:
    print(str(e))
    sys.exit(1)
    
from PyQt6.QtCore import Qt, pyqtSignal
from PyQt6.QtWidgets import (
    QApplication,
    QMainWindow,
    QWidget,
    QVBoxLayout,
    QGridLayout,
    QPushButton,
    QLabel,
)
from PyQt6.QtNetwork import QTcpSocket

class Packet(object):
    def __init__(self,
                 Keys: list[str],
                 Values: list[str]) -> None:
        self.keys = Keys
        self.values = Values

class Board(QWidget):
    tileSelected = pyqtSignal(int, int)
    
    def __init__(self, parent: QWidget | None) -> None:
        super().__init__(parent)
        
        self.glayout = QGridLayout(self)
    
        for row in range(3):
            for col in range(3):
                tile = QPushButton("", self)
                self.glayout.addWidget(tile, row, col, Qt.AlignmentFlag.AlignCenter)
                tile.pressed.connect(lambda x=row, y=col: self.tileSelected.emit(x, y))

class Client(QWidget):
    def __init__(self, parent: QWidget | None) -> None:
        super().__init__(parent)
        
        self.connectButton = QPushButton("Connect", self)
        self.connectButton.pressed.connect(self._connectPressed)
        
        self.board = Board(self)
        self.board.tileSelected.connect(self._tileSelected)
        
        self.vlayout = QVBoxLayout(self)
        self.vlayout.addWidget(self.connectButton)
        self.vlayout.addWidget(self.board)
        
        self.conn = QTcpSocket(self)
        self.conn.connected.connect(self._tcpConnected)
        self.conn.readyRead.connect(self._tcpRead)
        
    def _tcpConnected(self):
        print("Connected:" + self.conn.peerAddress().toString())
        print("Port:" + str(self.conn.peerPort()))
    
    def _tcpRead(self):
        data = self.conn.readAll().data().decode()
        j = json.loads(data)
        packet = Packet(**j)
        
        print(data)
        
    def _connectPressed(self):
        self.conn.connectToHost("127.0.0.1", 1200)
        
    def _tileSelected(self, x: int, y: int):
        print(x, y)
    
class MainWindow(QMainWindow):
    def __init__(self, parent: QWidget | None) -> None:
        super().__init__(parent)
        
        client = Client(self)
        self.setCentralWidget(client)

if __name__ == "__main__":
    app = QApplication(sys.argv)
    app.setApplicationName("Tic-Tac-Toe")
    
    mainwindow = MainWindow(None)
    mainwindow.show()
    
    app.exec()
