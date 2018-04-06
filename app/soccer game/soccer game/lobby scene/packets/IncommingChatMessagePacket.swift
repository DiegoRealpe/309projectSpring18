//
//  IncommingChatMessagePacket.swift
//  soccer game
//
//  Created by rtoepfer on 4/5/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

class IncommingChatMessagePacket {
    var playerNumber : Int
    var message : String
    
    init(data : [UInt8]){
        self.playerNumber = Int(data[1])
        
        let messageBytes = Array(data[2..<402])
        self.message = String(bytes: messageBytes, encoding: .utf8)!
        self.message = self.message.replacingOccurrences(of: "\0", with: "")
        
    }
    
}
