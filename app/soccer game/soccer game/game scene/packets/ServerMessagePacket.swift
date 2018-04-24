//
//  ServerMessagePacket.swift
//  soccer game
//
//  Created by rtoepfer on 4/23/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

class ServerMessagePacket {
    
    var message : String
    
    init(data: [UInt8] ){
        let messageBytes = Array(data[1..<81])
        self.message = String(bytes: messageBytes, encoding: .utf8)!
        self.message = self.message.replacingOccurrences(of: "\0", with: "")
        
    }
    
    
}
