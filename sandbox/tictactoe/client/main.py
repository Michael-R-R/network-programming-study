#!/usr/bin/python

import sys, json, types
from enum import Enum

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
    QLayout,
    QSizePolicy,
    QSpacerItem,
    QPushButton,
    QLabel,
)
from PyQt6.QtNetwork import QTcpSocket

class Message(str, Enum):
    BOARD_UPDATE = "10"
    PLAYER_STATE = "11"
    BANNER = "12"

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
        self.glayout.setHorizontalSpacing(0)
        self.glayout.setVerticalSpacing(0)
        self.glayout.setContentsMargins(0, 0, 0, 0)
    
        for row in range(3):
            for col in range(3):
                tile = QPushButton("", self)
                tile.setFixedSize(40, 40)
                self.glayout.addWidget(tile, row+1, col+1, Qt.AlignmentFlag.AlignCenter)
                
                tile.pressed.connect(lambda x=row, y=col: self.tileSelected.emit(x, y))
                
        self.glayout.addItem(QSpacerItem(20, 40, QSizePolicy.Policy.Minimum, QSizePolicy.Policy.Expanding), 0, 0, 3 + 2)
        self.glayout.addItem(QSpacerItem(20, 40, QSizePolicy.Policy.Minimum, QSizePolicy.Policy.Expanding), 3 + 1, 0, 1, 3 + 2)
        
        self.glayout.addItem(QSpacerItem(20, 40, QSizePolicy.Policy.Expanding, QSizePolicy.Policy.Minimum), 1, 0, 3, 1)
        self.glayout.addItem(QSpacerItem(20, 40, QSizePolicy.Policy.Expanding, QSizePolicy.Policy.Minimum), 1, 3 + 1, 3, 1)
                
    def update(self, text: str, row: int, col: int):
        item = self.glayout.itemAtPosition(row+1, col+1)
        button: QPushButton = item.widget()
        
        button.setText(text)

class Client(QWidget):
    def __init__(self, parent: QWidget | None) -> None:
        super().__init__(parent)
        
        self.playState = False
        
        self.connectButton = QPushButton("Connect", self)
        self.connectButton.pressed.connect(self._connectPressed)
        
        self.board = Board(self)
        self.board.tileSelected.connect(self._tileSelected)
        
        self.banner = QLabel(self)
        
        self.vlayout = QVBoxLayout(self)
        self.vlayout.addWidget(self.connectButton)
        self.vlayout.addWidget(self.board)
        self.vlayout.addWidget(self.banner)
        
        self.conn = QTcpSocket(self)
        self.conn.connected.connect(self._tcpConnected)
        self.conn.readyRead.connect(self._tcpRead)
        
    def _tcpConnected(self):
        self.connectButton.setEnabled(False)
        print("Connected:" + self.conn.peerAddress().toString())
        print("Port:" + str(self.conn.peerPort()))
    
    def _tcpRead(self):
        data = self.conn.readAll().data().decode()
        j = json.loads(data)
        packet = Packet(**j)
        
        # TODO parse packet values
        for i in range(len(packet.keys)):
            key = packet.keys[i]
            value = packet.values[i]
            
            match key:
                case Message.BOARD_UPDATE:
                    v = value.split(",")
                    self.board.update(v[0], int(v[1]), int(v[2]))
                case Message.PLAYER_STATE:
                    self.playState = True if (value == "1") else False
                case Message.BANNER:
                    self.banner.setText(value)
        
    def _connectPressed(self):
        self.conn.connectToHost("127.0.0.1", 1200)
        
    def _tileSelected(self, x: int, y: int):
        if self.playState == False:
            return
        
        self.playState = False
        
        keys: list[str] = [Message.BOARD_UPDATE]
        values: list[str] = ["{0},{1}".format(x, y)]
        packet = Packet(keys, values)
        
        jpckt = json.dumps(vars(packet))
        
        buf = bytearray()
        buf.extend(map(ord, jpckt))
        
        self.conn.write(buf)
    
class MainWindow(QMainWindow):
    def __init__(self, parent: QWidget | None) -> None:
        super().__init__(parent)
        
        client = Client(self)
        self.setCentralWidget(client)
        self.setFixedSize(200, 200)

if __name__ == "__main__":
    app = QApplication(sys.argv)
    app.setApplicationName("Tic-Tac-Toe")
    
    mainwindow = MainWindow(None)
    mainwindow.show()
    
    app.exec()
