//
//  PacketByteQueue.swift
//  soccer game
//
//  Created by rtoepfer on 2/16/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

class PacketByteQueue{
    
    var size = 0
    
    //oriented such that pop's happen from the first
    private var first : Node?
    private var last : Node?
    
    func deque(ammount : Int) -> [UInt8]{
        guard ammount <= size else{
            return []
        }
        
        var rtn = [UInt8]()
        rtn.reserveCapacity(ammount)
        
        for _ in 0..<ammount{
            rtn.append(self.deque())
        }
        self.size -= ammount
        
        return rtn
    }
    
    func enque(data : [UInt8]){
        self.size += data.count
        for value in data{
            enque(data : value)
        }
    }
    
    //does not modify size
    private func deque() -> UInt8{
        let rtn = first!.value
        first = first?.next //nil if next does not exist
        return rtn
    }
    
    //does not modify size
    private func enque(data : UInt8){
        let newNode = Node(data)
        last?.next = newNode
        last = newNode
        
        if self.first == nil{
            self.first = newNode
        }
    }
    
    class Node{
        var value : UInt8
        var next : Node?
        
        init(_ val: UInt8){
            self.value = val
        }
    }
    
}
