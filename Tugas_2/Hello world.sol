pragma solidity ^0.5.0;

//Mencoba compile file solidity di web remix.ethereum.org

contract halloWorld{
    uint256 totalCoin;

    //menambahkan saldo
    function addCoin(uint256 nCoin) public {
        totalCoin+=nCoin;
    }

    //melihat saldo
    function viewTotalCoin()public view returns(uint){
        return totalCoin;
    }  
}
