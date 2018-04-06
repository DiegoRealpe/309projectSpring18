//
//  OutgoingChatMessagePacket.swift
//  soccer game
//
//  Created by rtoepfer on 4/5/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

class OutgoingChatMessagePacket : SendablePacket {
    var message : String
    
    
    init(_ message : String){
        self.message = message
    }
    
    func toByteArray() -> [UInt8] {
        var array = [UInt8]()
        array.reserveCapacity(401)
        
        array.append(202)
        
        let utf8Array : [UInt8] = Array(message.utf8)
        array.append(contentsOf: utf8Array)
        array.append(contentsOf: Array(repeating: 0, count: 400-utf8Array.count))
        
        return array
    }
}
