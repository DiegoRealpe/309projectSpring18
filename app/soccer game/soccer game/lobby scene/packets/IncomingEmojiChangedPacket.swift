//
//  IncomingEmojiChangedPacket.swift
//  soccer game
//
//  Created by rtoepfer on 4/12/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

class IncomingEmojiChangedPacket {
    
    var playerNumber : Int
    var emoji : String
    
    init(_ data: [UInt8]) {
        self.playerNumber = Int(data[1])
        
        let messageBytes = Array(data[2..<25])
        self.emoji = String(bytes: messageBytes, encoding: .utf8)!
        self.emoji = self.emoji.replacingOccurrences(of: "\0", with: "")
    }
}
